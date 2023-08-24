# LogValuer

LogValuer is a tool to automate the creation of methods that satisfy the `slog.LogValuer` interface.

The implementation heavily references the `fmt.Stringer` interface.

## Usage

```go
//go:generate go run ../../cmd/logvaluer/logvaluer.go -type=Foo

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
```

Then, `go generate` creates a file which defines `LogValue` method for implementing `slog.LogValuer`.

```go
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
// 2023/08/24 22:14:57 INFO hello, world Foo.Str=hoge Foo.pstr=private Foo.Passwd=MASK Foo.Flo=3.14 Foo.M=1 Foo.N=20 Foo.L=300 Foo.Flag=true Foo.Date=2023-08-24T12:00:00.000+09:00 Foo.Dur=3s
```

## TODO

- Add support for type alias
- Add support for select multiple types at once
- Add support for exclude private struct field
