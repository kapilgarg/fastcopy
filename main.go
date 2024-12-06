package main

import (
	"fastcopy/copy/copy"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: fastcopy <source_file_path> <destination_directory> [<number_of_workers>]")
		os.Exit(1)
	}
	srcFilePath := os.Args[1]
	destDir := os.Args[2]
	numWorkers := 4
	if len(os.Args) == 4 {
		workers, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid number of workers. Provide a valid integer.")
			os.Exit(1)
		}
		numWorkers = workers
	}

	if err := copy.CopyFileFast(srcFilePath, destDir, numWorkers); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
