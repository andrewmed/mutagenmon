package extstat

import (
	"os"
	"syscall"
	"time"
)

func New(fi os.FileInfo) *ExtraStat {
	dateAttribute := fi.Sys().(*syscall.Win32FileAttributeData)
	return &ExtraStat{
		AccessTime: time.Unix(0, dateAttribute.LastAccessTime.Nanoseconds()),
		ModTime:    fi.ModTime(),
		ChangeTime: fi.ModTime(),
		BirthTime:  time.Unix(0, dateAttribute.CreationTime.Nanoseconds()),
	}
}
