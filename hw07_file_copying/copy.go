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

	//if fromPath == toPath {
	//	return ErrFilePathEqual
	//}

	src, err := os.Open(fromPath)

	if err != nil {
		// оборачиваем исходную ошибку. пример:  *os.PathError,
		// ErrOpenFile
		return fmt.Errorf("failed open file %w", err)
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			// просто выводим в консоль ошибку, для диагностики
			//
			fmt.Printf("failed close file %s", err)
		}
	}(src)

	//io.LimitedReader(src)

	// fileinfo
	fileInfo, err := src.Stat()
	if err != nil {
		//ErrGetFileInfo
		return ErrGetFileInfo
	}

	_, err = os.Stat(toPath)
	if err == nil {
		return ErrDestinationExistsFile
		//	return ErrSourceDestinationSameFile
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

	//coursor, err = src.Seek(offset, io.SeekStart)
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		//ErrCreateFile
		return fmt.Errorf("moving coursor was failed %w", err)
	}

	//reader := SetProgressBar(limit, src)
	//	_, err = io.Copy(dst)

	destination, err := os.Create(toPath)
	if err != nil {
		// ErrSeek
		return fmt.Errorf("moving coursor was failed %w", err)
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
		return fmt.Errorf("coping file was failed %w", err)
	}

	return nil
}

func CopyByPartial(dest *os.File, reader *pb.Reader, sizeBuffer int) error {

	if sizeBuffer < 1 {
		sizeBuffer = 1024
	}
	buf := make([]byte, sizeBuffer)

	//for item:=range buf {
	for {
		time.Sleep(1_00 * time.Millisecond)
		// func (r *Reader) Read(p []byte) (n int, err error) {
		countRead, errReadbuffer := reader.Read(buf) // двигаемся лимитировано через прогрессбар загруженный лимитировнный буфер
		if countRead > 0 {
			//countWrite, err := dest.Write(buf[:countRead])
			_, err := dest.Write(buf[:countRead]) // возврат count не нужно, записываем все что вычитали
			if err != nil {
				return fmt.Errorf("long copy error - write to buffer: %w", err)
			}
		}

		//
		// EOF is the error returned by Read when no more input is available.
		//var EOF = errors.New("EOF")
		if io.EOF == errReadbuffer { // мы все вычитали
			break
		}
		if errReadbuffer != nil {
			return fmt.Errorf("long copy error - read : %w", errReadbuffer)
		}

	}
	return nil
}

// TODO: начать с теста прогрессбара.
// https://github.com/cheggaaa/pb
// из примера: 37158 / 100000 [---------------->_______________________________] 37.16% 916 p/s
func SetProgressBar(limit int64, reader io.Reader) (*pb.Reader, *pb.ProgressBar) {
	bar := pb.Full.Start64(limit)

	// bar will format numbers as bytes (B, KiB, MiB, etc)
	bar.Set(pb.Bytes, true)

	//defer bar.Finish() // просто закрываем

	proxyReader := bar.NewProxyReader(io.LimitReader(reader, limit))
	return proxyReader, bar
}
