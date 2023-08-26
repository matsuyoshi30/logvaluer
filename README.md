# LogValuer

LogValuer is a tool to automate the creation of methods that satisfy the `slog.LogValuer` interface.

The implementation heavily references the `fmt.Stringer` interface.

## Usage

Execute `go get github.com/matsuyoshi30/logvaluer/cmd/logvaluer` and add `go:generate` comment. 


```go
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
```

Then, execute `go generate` creates a file which defines `LogValue` method for implementing `slog.LogValuer`.

```go
func (f Foo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Str", f.Str),
		slog.String("pstr", f.pstr),
		slog.String("Passwd", "MASK"),
		slog.Float64("Flo", f.Flo),
		slog.Int("M", f.M),
		slog.Int64("N", f.N),
		slog.Any("L", f.L),
		slog.Bool("Flag", f.Flag),
		slog.Time("Date", f.Date),
		slog.Duration("Dur", f.Dur),
	)
}
```

So you can pass the instance of this type to `slog` logging function.

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

- Add test
- Add support for type alias
- Add support for select multiple types at once
- Add support for exclude private struct field
