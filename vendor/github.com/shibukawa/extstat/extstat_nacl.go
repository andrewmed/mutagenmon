package extstat

import (
	"os"
	"syscall"
	"time"
)

func timespecToTime(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func New(fi os.FileInfo) *ExtraStat {
	osStat := fi.Sys().(*syscall.Stat_t)
	return &ExtraStat{
		AccessTime: timespecToTime(osStat.Atime, osStat.AtimeNsec),
		ModTime:    fi.ModTime(),
		ChangeTime: timespecToTime(osStat.Ctime, osStat.CtimeNsec),
		BirthTime:  fi.ModTime(),
	}
}
