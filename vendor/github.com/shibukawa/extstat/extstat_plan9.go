package extstat

import (
	"os"
	"time"
)

func New(fi os.FileInfo) *ExtraStat {
	return &ExtraStat{
		AccessTime: time.Unix(int64(fi.Sys().(*syscall.Dir).Atime), 0),
		ModTime:    fi.ModTime(),
		ChangeTime: fi.ModTime(),
		BirthTime:  fi.ModTime(),
	}
}
