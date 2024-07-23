package rds

import (
	"io"
	"strings"

	"github.com/common-fate/clio"
)

// DebugWriter is an io.Writer that writes messages using clio.Debug.
type DebugWriter struct{}

// Write implements the io.Writer interface for DebugWriter.
func (dw DebugWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	clio.Debug(message)
	return len(p), nil
}

type NotifyingWriter struct {
	writer   io.Writer
	phrase   string
	notifyCh chan struct{}
}

func NewNotifyingWriter(writer io.Writer, phrase string, notifyCh chan struct{}) *NotifyingWriter {
	return &NotifyingWriter{
		writer:   writer,
		phrase:   phrase,
		notifyCh: notifyCh,
	}
}

func (nw *NotifyingWriter) Write(p []byte) (n int, err error) {
	sentence := string(p)
	_ = sentence
	// Check if the phrase is in the input
	if strings.Contains(sentence, nw.phrase) {
		go func() { nw.notifyCh <- struct{}{} }()
	}
	// Write to the underlying writer
	return nw.writer.Write(p)
}
