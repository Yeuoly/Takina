package helper

import (
	"io"
	"math/rand"
	"os"
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
