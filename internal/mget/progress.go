package mget

import (
	"io"

	"github.com/schollz/progressbar/v3"
)

func ProgressWriter(totalBytes int64, description string) io.Writer {
	var bar *progressbar.ProgressBar

	if totalBytes > 0 {
		// Progress bar with known total
		bar = progressbar.DefaultBytes(totalBytes, description)
	} else {
		// Progress bar with unknown total
		bar = progressbar.DefaultBytes(-1, description)
	}

	return bar
}

func ProgressWriterUnknown(description string) io.Writer {
	return ProgressWriter(-1, description)
}
