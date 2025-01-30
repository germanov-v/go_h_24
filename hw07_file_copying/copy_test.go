package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.
}

func TestSetProgressBar(t *testing.T) {
	limit := int64(50)
	data := bytes.Repeat([]byte{1}, int(limit))
	reader := bytes.NewReader(data)

	progressReader := SetProgressBar(limit, reader)

	// Читаем данные из прогресс-бара
	readData, err := io.ReadAll(progressReader)
	//progressReader.Finish()
	//time.After(5 * time.Millisecond)

	if err != nil {
		t.Fatalf("Error reading: %v", err)
	}

	assert.Equal(t, data, readData, "Data must be equak")
}
