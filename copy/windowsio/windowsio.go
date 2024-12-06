package windowsio

import (
	"fastcopy/copy/progress"
	"os"

	"golang.org/x/sys/windows"
)

func CopyChunk(srcFile, destFile *os.File, offset, size int64, buffer []byte, tracker *progress.Tracker) {
	remaining := size
	for remaining > 0 {
		toRead := int64(len(buffer))
		if remaining < toRead {
			toRead = remaining
		}

		n, _ := pread(srcFile, buffer[:toRead], offset)
		pwrite(destFile, buffer[:n], offset)
		offset += int64(n)
		remaining -= int64(n)
		tracker.AddCopied(int64(n))
	}
}

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
