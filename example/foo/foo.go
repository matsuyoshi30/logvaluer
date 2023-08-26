package main

import (
	"log/slog"
	"time"
)

//go:generate go run github.com/matsuyoshi30/logvaluer/cmd/logvaluer -type=Foo

type Foo struct {
	Str    string
	pstr   string
	Passwd string `mask:"true"`
	Flo    float64

	M int
	N int64
	L uint64

	Flag bool

	Date time.Time
	Dur  time.Duration
}

func main() {
	foo := Foo{
		Str:    "hoge",
		pstr:   "private",
		Passwd: "",
		Flo:    3.14,
		M:      1,
		N:      20,
		L:      300,
		Flag:   true,
		Date:   time.Date(2023, time.August, 24, 12, 0, 0, 0, time.Local),
		Dur:    3 * time.Second,
	}
	slog.Info("hello, world", "Foo", foo)
}
