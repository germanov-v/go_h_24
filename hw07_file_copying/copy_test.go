package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestCopy(t *testing.T) {
	cleanDirTmp(t)
	// Place your code here.
	t.Run("with limit and offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "tmp/1.txt", 0, 10)
		assert.NoError(t, err)

		_, err = os.Stat("tmp/1.txt")
		assert.NoError(t, err)
	})

	t.Run("destination file exists", func(t *testing.T) {
		_, _ = os.Create("tmp/2.txt")
		err := Copy("testdata/input.txt", "tmp/2.txt", 0, 0)
		assert.ErrorIs(t, err, ErrDestinationExistsFile)
	})

	t.Run("null source", func(t *testing.T) {
		err := Copy("testdata/null.txt", "tmp/3.txt", 0, 0)
		assert.ErrorIs(t, err, ErrOpenFile)
	})

	t.Run("/dev/null", func(t *testing.T) {
		err := Copy("/dev/null", "tmp/4.txt", 0, 0)
		assert.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("error  offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "tmp/5.txt", 100000, 0)
		assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("without offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "tmp/6.txt", 0, 10000)
		assert.NoError(t, err)

		_, err = os.Stat("tmp/6.txt")
		assert.NoError(t, err)
	})

	t.Run("with offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "tmp/7.txt", 100, 0)
		assert.NoError(t, err)

		_, err = os.Stat("tmp/7.txt")
		assert.NoError(t, err)
	})
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

func cleanDirTmp(t *testing.T) {
	err := os.RemoveAll("tmp")
	assert.NoError(t, err)

	err = os.Mkdir("tmp", 0755)
	assert.NoError(t, err)
}
