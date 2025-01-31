package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	// Place your code here.
}

// calling: go test -v -run ^TestSetProgressBar$
func TestSetProgressBar(t *testing.T) {
	limit := int64(50)
	data := bytes.Repeat([]byte{1}, int(limit))
	reader := bytes.NewReader(data)

	progressReader, bar := SetProgressBar(limit, reader)

	bar.Finish()
	readData, err := io.ReadAll(progressReader)
	//progressReader.Finish()
	time.After(5 * time.Millisecond)

	if err != nil {
		t.Fatalf("Error reading: %v", err)
	}

	assert.Equal(t, data, readData, "Data must be equal")
}
