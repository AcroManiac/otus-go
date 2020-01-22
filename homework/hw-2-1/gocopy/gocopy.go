package gocopy

import (
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const copySize = 50 * 1024 // 50 KiBytes

func Copy(from string, to string, limit int, offset int) error {
	// Open reader
	file, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("could not open file for reading: %s", err.Error())
	}
	defer file.Close()

	// Evaluate bytes to write
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %s", err.Error())
	}
	bytesToWrite := info.Size()

	// Set limit of bytes to read
	input := io.Reader(file)
	if limit != -1 {
		input = io.LimitReader(file, int64(limit))
		bytesToWrite = int64(limit)
	}

	// Set offset position for reading
	if offset > 0 {
		offset64 := int64(offset)
		if pos, err := file.Seek(offset64, 0); err != nil || pos != offset64 {
			return fmt.Errorf("could not set offset for reading")
		}
		if offset64+bytesToWrite >= info.Size() {
			bytesToWrite = info.Size() - offset64
		}
	}

	// Create and start console progress bar
	bar := pb.StartNew(int(bytesToWrite))
	bar.SetWriter(os.Stdout)

	// Create writer
	output, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %s", err.Error())
	}
	defer output.Close()

	// Copy data
	for totalWritten := int64(0); totalWritten < bytesToWrite; {
		written, err := io.CopyN(output, input, copySize)
		if err != nil {
			if err == io.EOF {
				bar.Add64(written)
				break
			}
			return fmt.Errorf("error while copying the file: %s", err.Error())
		}
		totalWritten += written
		bar.Add64(written)
	}

	bar.Finish()
	return nil
}
