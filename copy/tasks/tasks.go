package tasks

import (
	"fastcopy/copy/progress"
	"fastcopy/copy/windowsio"

	"os"
)

type Task struct {
	Offset int64
	Size   int64
}

func CreateTaskQueue(fileSize, chunkSize int64, numWorkers int) chan Task {
	tasks := make(chan Task, numWorkers)
	for i := 0; i < numWorkers; i++ {
		offset := int64(i) * chunkSize
		size := chunkSize
		if i == numWorkers-1 {
			size = fileSize - offset
		}
		tasks <- Task{Offset: offset, Size: size}
	}
	close(tasks)
	return tasks
}

func ProcessTasks(taskQueue chan Task, srcFile, destFile *os.File, bufferSize int, progressTracker *progress.Tracker) {
	buffer := make([]byte, bufferSize)
	for task := range taskQueue {
		windowsio.CopyChunk(srcFile, destFile, task.Offset, task.Size, buffer, progressTracker)
	}
}
