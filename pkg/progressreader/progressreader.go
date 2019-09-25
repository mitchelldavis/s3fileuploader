package progressreader

import (
    "os"
    "sync/atomic"
)

type ProgressReader struct {
	File   *os.File
	Size int64
    Progress func(int64, int64)
	read int64
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	return r.File.Read(p)
}

func (r *ProgressReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.File.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&r.read, int64(n))

    if r.Progress != nil { r.Progress(r.Size, r.read) }

	return n, err
}

func (r *ProgressReader) Seek(offset int64, whence int) (int64, error) {
	return r.File.Seek(offset, whence)
}
