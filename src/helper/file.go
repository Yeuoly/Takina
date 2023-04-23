package helper

import (
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	root_temp_path = "/tmp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateTempFile() (string, io.ReadWriteCloser, error) {
	file, err := os.Create(
		root_temp_path +
			"/tmp-" + strconv.FormatInt(time.Now().Unix(), 10) +
			"-" + strconv.Itoa(rand.Int()))
	if err != nil {
		return "", nil, err
	}

	return file.Name(), file, nil
}

func CreateTempDir() (string, func(), error) {
	dir := root_temp_path +
		"/tmp-" + strconv.FormatInt(time.Now().Unix(), 10) +
		"-" + strconv.Itoa(rand.Int())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", nil, err
	}

	return dir, func() {
		// remove the temp dir
		os.RemoveAll(dir)
	}, nil
}

func WriteToFile(dir string, filename string, reader io.Reader) error {
	// create dir
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// create file
	file, err := os.Create(path.Join(dir, filename))
	if err != nil {
		return err
	}

	// write to file
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	// close
	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
