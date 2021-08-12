package extstat

import (
	"os"
	"syscall"
	"time"
)

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func New(fi os.FileInfo) *ExtraStat {
	osStat := fi.Sys().(*syscall.Stat_t)
	return &ExtraStat{
		AccessTime: timespecToTime(osStat.Atimespec),
		ModTime:    fi.ModTime(),
		ChangeTime: timespecToTime(osStat.Ctimespec),
		BirthTime:  timespecToTime(osStat.Birthtimespec),
	}
}
