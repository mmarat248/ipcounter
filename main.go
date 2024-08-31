package main

import (
	"awesomeProject/ipcounter"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	hyperLogLogType     = "hyperloglog"
	hyperLogLogPlusType = "hyperloglogplus"
	bitmapType          = "bitmap"
	setType             = "set"
)

func main() {
	// Define command-line flags
	filePath := flag.String("file", "", "Path to the file containing IP addresses")
	counterType := flag.String("counter", hyperLogLogPlusType, "Type of counter to use (hyperloglogplus or bitmap)")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
		log.Println("Please provide a file path using the -file flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create the counter based on the selected type
	var counter *ipcounter.IPCounter
	var err error
	switch *counterType {
	case hyperLogLogType:
		counter, err = createHyperLogLogCounter()
	case hyperLogLogPlusType:
		counter, err = createHyperLogLogPlusCounter()
	case bitmapType:
		counter, err = createBitmapCounter()
	case setType:
		counter, err = createSetCounter()
	default:
		log.Fatalf("Invalid counter type: %s", *counterType)
	}
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}

	start := time.Now()

	// Count IP addresses from the file
	count, err := counter.CountIPFromFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to count IP addresses from file %s: %v", *filePath, err)
	}

	elapsed := time.Since(start)

	fmt.Printf("%s count: %d\n", *counterType, count)
	fmt.Printf("Time elapsed: %v\n", elapsed)
}

func createHyperLogLogCounter() (*ipcounter.IPCounter, error) {
	hyperloglog, err := ipcounter.NewHyperLogLog(14)
	if err != nil {
		return nil, fmt.Errorf("failed to create HyperLogLogPlus: %v", err)
	}
	return ipcounter.NewIPCounter(hyperloglog, true, true), nil
}

func createHyperLogLogPlusCounter() (*ipcounter.IPCounter, error) {
	hyperloglogplus, err := ipcounter.NewHyperLogLogPlusBitMap(14)
	if err != nil {
		return nil, fmt.Errorf("failed to create HyperLogLogPlus: %v", err)
	}
	return ipcounter.NewIPCounter(hyperloglogplus, true, true), nil
}

func createBitmapCounter() (*ipcounter.IPCounter, error) {
	bitMap, err := ipcounter.NewIPBitMap()
	if err != nil {
		return nil, fmt.Errorf("failed to create BitmapCounter: %v", err)
	}
	return ipcounter.NewIPCounter(bitMap, true, false), nil
}

func createSetCounter() (*ipcounter.IPCounter, error) {
	s, _ := ipcounter.NewSet()
	return ipcounter.NewIPCounter(s, true, false), nil
}
