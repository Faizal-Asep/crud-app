package function

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateWriter struct {
	lock         sync.Mutex
	filename     string
	fp           *os.File
	maxAge       time.Duration
	rotationTime time.Duration
}

func NewLoger(filename string, rotationTime, maxAge time.Duration) *RotateWriter {
	w := &RotateWriter{filename: filename, rotationTime: rotationTime, maxAge: maxAge}
	err := w.Initialize()
	if err != nil {
		return nil
	}
	return w
}

func (w *RotateWriter) Initialize() (err error) {

	err = w.Rotate()
	if err != nil {
		return
	}

	rotation := time.NewTicker(w.rotationTime)

	go func() {
		for range rotation.C {
			w.Rotate()
		}
	}()
	return
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}

	file, err := os.Stat(w.filename)
	if err == nil {
		if time.Since(file.ModTime()) > w.rotationTime {
			err = os.Rename(w.filename, w.filename+"."+time.Now().Format(time.RFC3339))
			if err != nil {
				return
			}
		}

	}
	w.fp, err = os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// w.fp, err = os.Create()
	w.Age()
	return
}

func (w *RotateWriter) Age() (err error) {

	path := filepath.Dir(w.filename)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	isOlderThanMaxAge := func(t time.Time) bool {
		return time.Since(t) > w.maxAge
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if isOlderThanMaxAge(file.ModTime()) {
				os.Remove(filepath.Join(path, file.Name()))

			}
		}
	}

	return
}
