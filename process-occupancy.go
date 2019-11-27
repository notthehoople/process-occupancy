package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Format Header:	  TimeStamp,Date,Time,UniqueDevices,WiredDevices,WirelessDevices,Office
// Example of format: 2019-08-20-23:40,2019-08-20,23:40,28,17,16,London
func processOccupancyData(londonData bool, debug bool) {
	var uniqueDevices, minUnique, maxUnique int = 0, 2000, 0
	var outputDate string
	var directoryName = "testdata/London/"

	// Setup at process start

	if debug {
		fmt.Println("Output something interesting to debug")
	}

	files, err := ioutil.ReadDir(directoryName)
	if err != nil {
		log.Fatal(err)
	}

	for _, currentFile := range files {
		if debug {
			fmt.Println(currentFile.Name())
		}

		// Open the file
		csvFile, err := os.Open(directoryName + currentFile.Name())
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
				outputDate = csvRecord[1]
			} else {
				// Look for smallest number of devices, and ignore the collection bug that reports randomly small numbers
				if uniqueDevices < minUnique && uniqueDevices > 5 {
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

		// Output for each file: Filename, Date, Minimum Unique Devices, Maximum Unique Devices
		fmt.Printf("%s,%s,%d,%d\n", currentFile.Name(), outputDate, minUnique, maxUnique)
		maxUnique = 0
		minUnique = 2000
	}
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
