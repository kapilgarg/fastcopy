package copy

import (
	"fastcopy/copy/progress"
	"fastcopy/copy/tasks"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const bufferSize = 4 * 1024 * 1024 // 4 MB

// CopyFileFast handles the high-performance copying of files.
func CopyFileFast(src, destDir string, numWorkers int) error {
	srcFile, srcFileSize, err := openSourceFile(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, destFile, err := createDestinationFile(src, destDir)
	if err != nil {
		return err
	}
	defer destFile.Close()

	chunkSize := srcFileSize / int64(numWorkers)
	taskQueue := tasks.CreateTaskQueue(srcFileSize, chunkSize, numWorkers)
	progressTracker := progress.NewTracker(srcFileSize)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tasks.ProcessTasks(taskQueue, srcFile, destFile, bufferSize, progressTracker)
		}()
	}

	progress.StartProgressTicker(progressTracker)
	wg.Wait()
	progress.StopProgressTicker(progressTracker)

	fmt.Println("\nCopy complete!")
	return nil
}

func openSourceFile(src string) (*os.File, int64, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open source file: %v", err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get source file info: %v", err)
	}
	return file, fileInfo.Size(), nil
}

func createDestinationFile(src, destDir string) (string, *os.File, error) {
	destPath := filepath.Join(destDir, filepath.Base(src))
	file, err := os.Create(destPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	return destPath, file, nil
}
