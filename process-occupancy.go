package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Read the text file passed in by name into a array of strings
// Returns the array as the first return variable
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func printStringArray(tempString []string) {
	// Loop through the array and print each line
	for i := 0; i < len(tempString); i++ {
		fmt.Println(tempString[i])
	}
}

// Format Header:	  TimeStamp,Date,Time,UniqueDevices,WiredDevices,WirelessDevices,Office
// Example of format: 2019-08-20-23:40,2019-08-20,23:40,28,17,16,London
func processOccupancyData(londonData bool, debug bool) {
	var uniqueDevices, minUnique, maxUnique int = 0, 2000, 0
	var fileName = "testdata/London/2019-08-20-Tuesday.txt"

	// Setup at process start

	if debug {
		fmt.Println("Output something interesting to debug")
	}

	// Open the file
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	reader := csv.NewReader(csvFile)
	// Read and ignore the first record
	reader.Read()

	// Iterate through the records
	for {
		// Read each record from csv
		csvRecord, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// We really care about field 3 - Unique Devices
		// We'll capture Date (1) and Time (2) as well for debug reasons
		uniqueDevices, _ = strconv.Atoi(csvRecord[3])
		if uniqueDevices > maxUnique {
			if debug {
				fmt.Println("Biggest seen")
			}
			maxUnique = uniqueDevices
		} else {
			// Look for smallest number of devices, and ignore the collection bug that reports "1" randomly
			if uniqueDevices < minUnique && uniqueDevices != 1 {
				if debug {
					fmt.Println("Smallest seen")
				}

				minUnique = uniqueDevices
			}
		}

		if debug {
			fmt.Printf("Date: %s Time: %s uniqueDevices: %d\n", csvRecord[1], csvRecord[2], uniqueDevices)
		}
	}

	fmt.Printf("Smallest: %d Biggest: %d\n", minUnique, maxUnique)
}

// Main routine
func main() {
	var debug bool = false
	var londonData bool = false

	flag.BoolVar(&debug, "debug", false, "Turn debug on")
	flag.BoolVar(&londonData, "london", false, "Process London data")

	flag.Parse()

	processOccupancyData(londonData, debug)
}
