package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.

	src, err := os.Open(fromPath)

	if err != nil {
		// оборачиваем исходную ошибку. пример:  *os.PathError,
		// ErrUnsupportedFile
		return fmt.Errorf("failed open file %w", err)
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			// просто выводим в консоль ошибку, для диагностики
			fmt.Printf("failed close file %s", err)
		}
	}(src)

	//io.LimitedReader(src)

	//reader := SetProgressBar(limit, src)

	//	_, err = io.Copy(dst)

	return nil
}

// TODO: начать с теста прогрессбара.
// https://github.com/cheggaaa/pb
// из примера: 37158 / 100000 [---------------->_______________________________] 37.16% 916 p/s
func SetProgressBar(limit int64, reader io.Reader) *pb.Reader {
	bar := pb.Full.Start64(limit)

	// bar will format numbers as bytes (B, KiB, MiB, etc)
	bar.Set(pb.Bytes, true)

	//defer bar.Finish() // просто закрываем

	proxyReader := bar.NewProxyReader(io.LimitReader(reader, limit))
	return proxyReader
}
