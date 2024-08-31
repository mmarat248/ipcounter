package ipcounter

import (
	"os"
	"testing"
)

// MockIPMap is a mock implementation of IPMap for testing
type MockIPMap struct {
	count uint64
	ips   map[uint32]struct{}
}

func NewMockIPMap() *MockIPMap {
	return &MockIPMap{
		ips: make(map[uint32]struct{}),
	}
}

func (m *MockIPMap) Add(ip uint32) {
	if _, exists := m.ips[ip]; !exists {
		m.ips[ip] = struct{}{}
		m.count++
	}
}

func (m *MockIPMap) Count() uint64 {
	return m.count
}

func TestCountIPFromFile(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected uint64
	}{
		{"Single IP", "192.168.0.1\n", 1},
		{"Multiple IPs", "192.168.0.1\n10.0.0.1\n192.168.0.1\n", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file
			tmpfile, err := os.CreateTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write content to the file
			if _, err := tmpfile.Write([]byte(tc.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Create counter and count IPs
			mockMap := NewMockIPMap()
			counter := NewIPCounter(mockMap, false, true)
			count, err := counter.CountIPFromFile(tmpfile.Name())

			// Check results
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if count != tc.expected {
				t.Errorf("Expected count %d, got %d", tc.expected, count)
			}
		})
	}
}

func TestParseIP(t *testing.T) {
	testCases := []struct {
		input    string
		expected uint32
	}{
		{"192.168.0.1", 3232235521},
		{"10.0.0.1", 167772161},
		{"172.16.0.1", 2886729729},
	}

	for _, tc := range testCases {
		result := parseIP([]byte(tc.input))
		if result != tc.expected {
			t.Errorf("For input %s, expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

func TestProcessChunk(t *testing.T) {
	mockMap := NewMockIPMap()
	counter := NewIPCounter(mockMap, false, true)

	chunk := []byte("192.168.0.1\n10.0.0.1\n192.168.0.1\n")
	err := counter.processChunk(chunk)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mockMap.Count() != 2 {
		t.Errorf("Expected 2 unique IPs, got %d", mockMap.Count())
	}
}

func TestAddIPBatch(t *testing.T) {
	mockMap := NewMockIPMap()
	counter := NewIPCounter(mockMap, false, true)

	ips := []uint32{3232235521, 167772161, 3232235521} // 192.168.0.1, 10.0.0.1, 192.168.0.1
	counter.addIPBatch(ips)

	if mockMap.Count() != 2 {
		t.Errorf("Expected 2 unique IPs, got %d", mockMap.Count())
	}
}
