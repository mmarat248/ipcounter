package ipcounter

import (
	"awesomeProject/ipcounter/utils/fnv1a"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
)

// IPBatchSize The number of IPs to process in a single batch
const IPBatchSize = 200

// IPCounter represents a structure for counting unique IP addresses
type IPCounter struct {
	ipMap       IPMap
	useParallel bool
	lock        sync.Mutex
	useHashFunc bool
}

func NewIPCounter(mp IPMap, useParallel, useHashFunc bool) *IPCounter {
	return &IPCounter{
		ipMap:       mp,
		useParallel: useParallel,
		useHashFunc: useHashFunc,
	}
}

// CountIPFromFile counts unique IPs from a file
func (counter *IPCounter) CountIPFromFile(fileName string) (uint64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file stat: %w", err)
	}
	fileSize := fileStat.Size()
	if fileSize <= 0 || fileSize != int64(int(fileSize)) {
		return 0, fmt.Errorf("wrong file size: %d", fileStat.Size())
	}

	fileChunkSize := 1073741824 // should % 4096
	offset := int64(0)
	for offset < fileStat.Size() {
		length := fileChunkSize
		if length > int(fileStat.Size()-offset) {
			length = int(fileStat.Size() - offset)
		}
		err = counter.ProcessFileChunk(file, offset, length)
		if err != nil {
			return 0, err
		}
		offset += int64(fileChunkSize)
	}
	return counter.ipMap.Count(), nil
}

func (counter *IPCounter) ProcessFileChunk(file *os.File, fileChunkOffset int64, fileChunkLength int) error {
	// Install a page fault handler, so that I/O errors against the
	// memory map (e.g., due to disk failure) don't cause us to
	// crash.
	prevPanicOnFault := debug.SetPanicOnFault(true)
	defer func() {
		debug.SetPanicOnFault(prevPanicOnFault)
		if recover() != nil {
			log.Print("Page fault occurred while reading from memory map")
		}
	}()

	data, err := syscall.Mmap(int(file.Fd()), fileChunkOffset, fileChunkLength, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("failed to mmap file: %v", err)
	}
	defer func() {
		if err := syscall.Munmap(data); err != nil {
			log.Printf("failed to munmap: %v", err)
		}
	}()

	// Process all data in a single chunk
	if !counter.useParallel {
		return counter.processChunk(data)
	}

	// Split data into chunkEndPositions for parallel processing
	chunkSize, chunkCount := getChunkSize(len(data))
	chunkEndPositions := make([]int, 0, chunkCount)
	offset := 0
	for offset < len(data) {
		offset += chunkSize
		if offset >= len(data) {
			chunkEndPositions = append(chunkEndPositions, len(data))
			break
		}

		// Find the next newline to ensure we don't split an IP address
		newlinePosition := bytes.IndexByte(data[offset:], '\n')
		if newlinePosition == -1 {
			chunkEndPositions = append(chunkEndPositions, len(data))
			break
		} else {
			offset += newlinePosition + 1
			chunkEndPositions = append(chunkEndPositions, offset)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(chunkEndPositions))
	errChan := make(chan error, chunkCount)

	chunkStart := 0
	for _, chunkEnd := range chunkEndPositions {
		go func(chunk []byte) {
			defer wg.Done()
			if chunkErr := counter.processChunk(chunk); chunkErr != nil {
				errChan <- chunkErr
			}
		}(data[chunkStart:chunkEnd])
		chunkStart = chunkEnd
	}

	wg.Wait()
	close(errChan)
	for err = range errChan {
		if err != nil {
			// TODO: merge errors
			return err
		}
	}
	return nil
}

// processChunk processes a chunk of data and counts IPs
func (counter *IPCounter) processChunk(data []byte) error {
	ipBatch := make([]uint32, 0, IPBatchSize)

	for len(data) > 0 {
		endOfLine := bytes.IndexByte(data, '\n')
		if endOfLine == -1 {
			endOfLine = len(data)
		}

		ip := parseIP(data[:endOfLine])
		ipBatch = append(ipBatch, ip)
		if len(ipBatch) == IPBatchSize {
			counter.addIPBatch(ipBatch)
			ipBatch = ipBatch[:0]
		}

		if endOfLine == len(data) {
			break
		}
		data = data[endOfLine+1:]
	}

	if len(ipBatch) > 0 {
		counter.addIPBatch(ipBatch)
	}
	return nil
}

// addIPBatch adds a batch of IPs to the counter
func (counter *IPCounter) addIPBatch(ips []uint32) {
	if counter.useParallel {
		counter.lock.Lock()
		defer counter.lock.Unlock()
	}
	for _, ip := range ips {
		if counter.useHashFunc {
			ip = fnv1a.HashUint32(ip)
		}
		counter.ipMap.Add(ip)
	}
}

// parseIP converts a byte slice to a uint32 IP representation
func parseIP(data []byte) uint32 {
	var ip uint32
	var octet uint8

	for _, b := range data {
		if b == '.' {
			ip = (ip << 8) | uint32(octet)
			octet = 0
			continue
		}
		octet = (octet * 10) + (b - '0')
	}
	return (ip << 8) | uint32(octet)
}

func getChunkSize(dataLength int) (int, int) {
	nChunks := runtime.NumCPU()
	chunkSize := dataLength / nChunks
	if chunkSize == 0 {
		chunkSize = dataLength
	}
	return chunkSize, nChunks
}
