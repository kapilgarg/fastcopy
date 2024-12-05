package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sys/windows"
)

const bufferSize = 4 * 1024 * 1024 // 4 MB

func CopyFileFast(src, destDir string, numWorkers int) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	srcFileSize := srcFileInfo.Size()

	destFilePath := filepath.Join(destDir, filepath.Base(src))
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	chunkSize := srcFileSize / int64(numWorkers)
	type Task struct {
		Offset int64
		Size   int64
	}

	tasks := make(chan Task, numWorkers)

	var totalCopied int64
	progressMutex := sync.Mutex{}

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		buffer := make([]byte, bufferSize)
		for task := range tasks {
			remaining := task.Size
			offset := task.Offset

			for remaining > 0 {
				toRead := int64(bufferSize)
				if remaining < toRead {
					toRead = remaining
				}

				n, err := pread(srcFile, buffer[:toRead], offset)
				if err != nil && err != io.EOF {
					fmt.Printf("\nError reading: %v\n", err)
					return
				}

				_, err = pwrite(destFile, buffer[:n], offset)
				if err != nil {
					fmt.Printf("\nError writing: %v\n", err)
					return
				}

				offset += int64(n)
				remaining -= int64(n)

				progressMutex.Lock()
				totalCopied += int64(n)
				progressMutex.Unlock()
			}
		}
	}

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Start a ticker to display progress and rate
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	done := make(chan struct{})

	go func() {
		var lastCopied int64
		for {
			select {
			case <-ticker.C:
				progressMutex.Lock()
				bytesCopied := totalCopied - lastCopied
				lastCopied = totalCopied
				rate := float64(bytesCopied) / float64(1) // Rate in bytes per second
				percent := int((totalCopied * 100) / srcFileSize)
				fmt.Printf("\rCopying... %d%% | %.2f MB/s", percent, rate/1024/1024)
				progressMutex.Unlock()
			case <-done:
				progressMutex.Lock()
				percent := int((totalCopied * 100) / srcFileSize)
				fmt.Printf("\rCopying... %d%% | %.2f MB/s", percent, 0.0) // Ensure final rate is shown as 0
				progressMutex.Unlock()
				return
			}
		}
	}()

	startTime := time.Now()
	// Distribute tasks
	for i := 0; i < numWorkers; i++ {
		offset := int64(i) * chunkSize
		size := chunkSize
		if i == numWorkers-1 {
			size = srcFileSize - offset // Last chunk may be larger
		}
		tasks <- Task{Offset: offset, Size: size}
	}
	close(tasks)

	wg.Wait()

	close(done)

	elapsedTime := time.Since(startTime)
	fmt.Printf("\nCopy complete! Time taken: %.2fs\n", elapsedTime.Seconds())
	return nil
}

// pread reads from the file at a specific offset using Windows system calls.
func pread(file *os.File, buffer []byte, offset int64) (int, error) {
	fileHandle := windows.Handle(file.Fd())
	var bytesRead uint32
	overlapped := windows.Overlapped{Offset: uint32(offset), OffsetHigh: uint32(offset >> 32)}

	err := windows.ReadFile(fileHandle, buffer, &bytesRead, &overlapped)
	if err != nil && err != windows.ERROR_HANDLE_EOF {
		return 0, err
	}
	return int(bytesRead), nil
}

// pwrite writes to the file at a specific offset using Windows system calls.
func pwrite(file *os.File, buffer []byte, offset int64) (int, error) {
	fileHandle := windows.Handle(file.Fd())
	var bytesWritten uint32
	overlapped := windows.Overlapped{Offset: uint32(offset), OffsetHigh: uint32(offset >> 32)}

	err := windows.WriteFile(fileHandle, buffer, &bytesWritten, &overlapped)
	if err != nil {
		return 0, err
	}
	return int(bytesWritten), nil
}

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

	if err := CopyFileFast(srcFilePath, destDir, numWorkers); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
