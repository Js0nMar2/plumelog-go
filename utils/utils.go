package utils

import (
	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
	"plumelog/log"
	"runtime"
	"strconv"
	"time"
)

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func WorkId() string {
	newSonyflake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := newSonyflake.NextID()
	if err != nil {
		log.Error.Println(err)
	}
	return strconv.Itoa(int(id))
}

func SameDay(beginDate, endDate int64) bool {
	y1, m1, d1 := time.UnixMilli(beginDate).Date()
	y2, m2, d2 := time.UnixMilli(endDate).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ParseDay(datetime int64, pattern string) string {
	return time.UnixMilli(datetime).Format(pattern)
}

func NextDay(datetime int64, days int, pattern string) string {
	return time.UnixMilli(datetime).AddDate(0, 0, days).Format(pattern)
}

func BetweenDays(beginDate, endDate int64) int {
	return int((endDate - beginDate) / (1000 * 60 * 60 * 24))
}
