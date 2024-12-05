package main

import (
	"bufio"         //for buffered I/O
	"encoding/json" //used to encode data into json format
	"fmt"
	"os"      // provides function to interact with operating system
	"regexp"  //provides support to regular expression
	"strings" // provides string manipulation function
	"sync"    //here it provides synchronisation primitives
)

type LogEntry struct {
	Timestamp string `json:"timestamp"` //the timestamp of the log entry
	Level     string `json:"level"`     //the log level("errore","info")
	Message   string `json:"message"`   //log meassage
}

// this pattern (logPattern) is used to extract log information from each line of the log file.
var logPattern = regexp.MustCompile(`^(?P<timestamp>\S+) \[(?P<level>[A-Z]+)\] (?P<message>.+)$`)

//This function takes a line from the log file and attempts to match it against the logPattern regular expression.

func parseLogLine(line string) (LogEntry, bool) {
	matches := logPattern.FindStringSubmatch(line)
	if matches == nil {
		return LogEntry{}, false
	}
	return LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, true
}

// This function processes a chunk of lines (a slice of strings).

func handleChunk(lines []string, results chan LogEntry, seen *sync.Map) {
	for _, line := range lines {
		if logEntry, valid := parseLogLine(line); valid {
			// Check for duplicates using a concurrent map
			entryKey := fmt.Sprintf("%s|%s|%s", logEntry.Timestamp, logEntry.Level, logEntry.Message)
			if _, exists := seen.LoadOrStore(entryKey, struct{}{}); !exists {
				results <- logEntry
			}
		}
	}
}

// This function reads the log file in chunks of a specified size
func readLogFileInChunks(filepath string, chunkSize int) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]string
	scanner := bufio.NewScanner(file)
	var chunk []string
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) >= chunkSize {
			chunks = append(chunks, chunk)
			chunk = nil
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks, scanner.Err()
}

// This new function checks the log entries for validity.

func validateLogEntries(logEntries []LogEntry) error {
	for _, entry := range logEntries {
		// Example validation: Check if log level is valid (ERROR, INFO, etc.)
		if entry.Level != "ERROR" && entry.Level != "INFO" && entry.Level != "WARN" {
			return fmt.Errorf("invalid log level: %s", entry.Level)
		}
	}
	return nil
}

// This function saves the log entries into a JSON file.

func saveLogEntriesToJSON(logEntries []LogEntry, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logEntries); err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}
	return nil
}

func main() {

	filepath := "C:/Users/khusboo.k/Desktop/Mygolang/task2/Day-4_log_file_500mb.log"
	outputFile := strings.Replace(filepath, ".log", ".json", 1)
	chunkSize := 1000

	// Read all file in chunks
	chunks, err := readLogFileInChunks(filepath, chunkSize)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	results := make(chan LogEntry, chunkSize)
	var wg sync.WaitGroup
	var seen sync.Map

	// Process chunks concurrently
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			handleChunk(chunk, results, &seen)
		}(chunk)
	}

	// Close results channel after all processing is done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var logEntries []LogEntry
	for entry := range results {
		logEntries = append(logEntries, entry)
	}

	// now saving results to JSON file
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logEntries); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

}
