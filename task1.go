package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// File to read and write to
const CSVFile = "C:/Users/khusboo.k/Desktop/Mygolang/go.mod/fixlets.csv"

// Structure to represent a CSV entry
type Entry struct {
	SiteID                string
	FxiletID              string
	Name                  string
	Criticality           string
	RelevantComputerCount int
}

// Helper function to read the CSV file
func readCSV() ([]Entry, error) {
	var entries []Entry

	file, err := os.Open(CSVFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	// Read each row and append to entries slice
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Convert RelevantComputerCount to integer
		relevantComputerCount, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, fmt.Errorf("invalid value for RelevantComputerCount: %v", err)
		}

		entry := Entry{
			SiteID:                record[0],
			FxiletID:              record[1],
			Name:                  record[2],
			Criticality:           record[3],
			RelevantComputerCount: relevantComputerCount,
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// Helper function to write the entries to the CSV file
func writeCSV(entries []Entry) error {
	file, err := os.Create(CSVFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	writer.Write([]string{"SiteID", "FxiletID", "Name", "Criticality", "RelevantComputerCount"})

	// Write the entries
	for _, entry := range entries {
		writer.Write([]string{
			entry.SiteID,
			entry.FxiletID,
			entry.Name,
			entry.Criticality,
			strconv.Itoa(entry.RelevantComputerCount),
		})
	}

	return nil
}

// List all entries in the CSV
func listEntries() {
	entries, err := readCSV()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	fmt.Println("Listing Entries:")
	for _, entry := range entries {
		fmt.Printf("SiteID: %s, FxiletID: %s, Name: %s, Criticality: %s, RelevantComputerCount: %d\n",
			entry.SiteID, entry.FxiletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
	}
}

// Query an entry by FxiletID
func queryEntry() {
	var fxiletID string
	fmt.Print("Enter the FxiletID to search for: ")
	fmt.Scanln(&fxiletID)

	entries, err := readCSV()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	found := false
	for _, entry := range entries {
		if entry.FxiletID == fxiletID {
			fmt.Printf("Found: SiteID: %s, FxiletID: %s, Name: %s, Criticality: %s, RelevantComputerCount: %d\n",
				entry.SiteID, entry.FxiletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("No entry found.")
	}
}

// Sort entries by Criticality
func sortEntries() {
	entries, err := readCSV()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Sorting entries by Criticality
	for i := 0; i < len(entries)-1; i++ {
		for j := 0; j < len(entries)-i-1; j++ {
			if strings.Compare(entries[j].Criticality, entries[j+1].Criticality) > 0 {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}

	// Display sorted entries
	fmt.Println("Sorted Entries by Criticality:")
	for _, entry := range entries {
		fmt.Printf("SiteID: %s, FxiletID: %s, Name: %s, Criticality: %s, RelevantComputerCount: %d\n",
			entry.SiteID, entry.FxiletID, entry.Name, entry.Criticality, entry.RelevantComputerCount)
	}

	// Save sorted entries back to CSV
	err = writeCSV(entries)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
	}
}

// Add a new entry to the CSV
func addEntry() {
	var siteID, fxiletID, name, criticality, relevantComputerCountStr string
	fmt.Print("Enter SiteID: ")
	fmt.Scanln(&siteID)
	fmt.Print("Enter FxiletID: ")
	fmt.Scanln(&fxiletID)
	fmt.Print("Enter Name: ")
	fmt.Scanln(&name)
	fmt.Print("Enter Criticality: ")
	fmt.Scanln(&criticality)
	fmt.Print("Enter RelevantComputerCount: ")
	fmt.Scanln(&relevantComputerCountStr)

	// Convert relevantComputerCount to integer
	relevantComputerCount, err := strconv.Atoi(relevantComputerCountStr)
	if err != nil {
		fmt.Println("Error converting RelevantComputerCount to integer:", err)
		return
	}

	entry := Entry{
		SiteID:                siteID,
		FxiletID:              fxiletID,
		Name:                  name,
		Criticality:           criticality,
		RelevantComputerCount: relevantComputerCount,
	}

	// Read existing entries
	entries, err := readCSV()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Add the new entry
	entries = append(entries, entry)

	// Write back to CSV
	err = writeCSV(entries)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return
	}

	fmt.Println("Entry added.")
}

// Delete an entry by FxiletID
func deleteEntry() {
	var fxiletID string
	fmt.Print("Enter the FxiletID of the entry to delete: ")
	fmt.Scanln(&fxiletID)

	entries, err := readCSV()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Filter out the entry to delete
	var newEntries []Entry
	for _, entry := range entries {
		if entry.FxiletID != fxiletID {
			newEntries = append(newEntries, entry)
		}
	}

	// Write back to CSV
	err = writeCSV(newEntries)
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return
	}

	fmt.Println("Entry deleted.")
}

// Main menu to interact with the user
func main() {
	for {
		fmt.Println("\nCSV Management Menu:")
		fmt.Println("1. List all entries")
		fmt.Println("2. Query an entry by FxiletID")
		fmt.Println("3. Sort entries by Criticality")
		fmt.Println("4. Add a new entry")
		fmt.Println("5. Delete an entry by FxiletID")
		fmt.Println("6. Exit")

		var choice int
		fmt.Print("Enter your choice (1-6): ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			listEntries()
		case 2:
			queryEntry()
		case 3:
			sortEntries()
		case 4:
			addEntry()
		case 5:
			deleteEntry()
		case 6:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
