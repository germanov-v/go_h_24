package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrOpenFile                  = errors.New("open file failed")
	ErrGetFileInfo               = errors.New("file info failed")
	ErrSourceDestinationSameFile = errors.New("source file destination are same file")
	ErrDestinationExistsFile     = errors.New("destination file exists")

	ErrSeekFile      = errors.New("seek failed")
	ErrCopyFile      = errors.New("copy file failed")
	ErrCreateFile    = errors.New("create file failed")
	ErrFilePathEqual = errors.New("file path equal")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	_, errDestination := os.Stat(toPath)
	if errDestination == nil {
		return ErrDestinationExistsFile
	}

	src, err := os.Open(fromPath)

	if err != nil {
		// оборачиваем исходную ошибку. пример:  *os.PathError,
		// ErrOpenFile
		return ErrOpenFile
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			// просто выводим в консоль ошибку, для диагностики
			//
			fmt.Printf("failed close file %s", err)
		}
	}(src)

	fileInfo, err := src.Stat()
	if err != nil {
		//ErrGetFileInfo
		return ErrGetFileInfo
	}

	// под ErrUnsupportedFile: не /dev/nukl, dir, slink, (socket?)
	if !fileInfo.Mode().IsRegular() {
		//  return fmt.Errorf("%s is not a regular file", fromPath)
		return ErrUnsupportedFile
	}

	size := fileInfo.Size()
	if size < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || size-offset < limit {
		// выход за границы файла - берем все
		limit = size - offset
	}
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		//ErrCreateFile
		return ErrSeekFile
	}

	destination, err := os.Create(toPath)
	if err != nil {
		// ErrSeek
		return ErrCreateFile
	}

	defer destination.Close()

	limitReader := io.LimitReader(src, limit)

	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	bar.SetRefreshRate(time.Millisecond)

	bar.SetWriter(os.Stdout)
	defer bar.Finish()
	//reader, _ := SetProgressBar(limit, limitReader)

	reader := bar.NewProxyReader(limitReader)
	time.Sleep(time.Millisecond)
	_, err = io.Copy(destination, reader)
	//err = CopyByPartial(destination, reader, 5)
	if err != nil && err != io.EOF {
		// ErrCopyFile
		return ErrCopyFile
	}

	return nil
}

func CopyByPartial(dest *os.File, reader *pb.Reader, sizeBuffer int) error {

	if sizeBuffer < 1 {
		sizeBuffer = 1024
	}
	buf := make([]byte, sizeBuffer)

	for {
		time.Sleep(1_00 * time.Millisecond)
		countRead, errReadbuffer := reader.Read(buf) // двигаемся лимитировано через прогрессбар загруженный лимитировнный буфер
		if countRead > 0 {
			_, err := dest.Write(buf[:countRead]) // возврат count не нужно, записываем все что вычитали
			if err != nil {
				return fmt.Errorf("long copy error - write to buffer: %w", err)
			}
		}

		// EOF is the error returned by Read when no more input is available.
		if io.EOF == errReadbuffer { // мы все вычитали
			break
		}
		if errReadbuffer != nil {
			return fmt.Errorf("long copy error - read : %w", errReadbuffer)
		}

	}
	return nil
}

// https://github.com/cheggaaa/pb
func SetProgressBar(limit int64, reader io.Reader) (*pb.Reader, *pb.ProgressBar) {
	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)

	proxyReader := bar.NewProxyReader(io.LimitReader(reader, limit))
	return proxyReader, bar
}
