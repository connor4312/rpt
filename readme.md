# rpt [![Build Status](https://travis-ci.org/connor4312/rpt.svg?branch=master)](https://travis-ci.org/connor4312/rpt) [![Coverage Status](https://coveralls.io/repos/connor4312/rpt/badge.svg?branch=master)](https://coveralls.io/r/connor4312/rpt?branch=master) [![godoc reference](https://godoc.org/github.com/connor4312/rpt?status.png)](https://godoc.org/github.com/connor4312/rpt)

RPT ("requests per time") is a general-purpose library for monitoring events over time interval. For example, it can be used to easily and quickly calculate the number of requests your app is hit by per minute.

See the godoc for further details.

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/connor4312/rpt"
    "net/http"
    "time"
)

func main() {
    rp := rpt.New(60, time.Second)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        rp.AddRequest()
        fmt.Fprintf(w, "Requests per minute: %d", rp.GetRPT())
    })

    http.ListenAndServe(":8080", nil)
}
```

## Benchmarks

 * The call to record a request takes under 30 nanoseconds.
 * The call to sum all requests in the interval takes about 175 nanoseconds.
 * The call to get a range of data, suitable for building histograms, takes about 900 nanoseconds.

```
➜  rpt  go test -benchmem -bench=.
PASS
BenchmarkAddRequest 50000000          28.8 ns/op           0 B/op          0 allocs/op
BenchmarkGetRtp     10000000           170 ns/op           0 B/op          0 allocs/op
BenchmarkGetRange    2000000           879 ns/op         704 B/op          1 allocs/op
ok      github.com/connor4312/worker/rpt 5.949s
```

## License

Copyright 2015 by Connor Peet. Distributed under the MIT license.
