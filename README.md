# BeamNG SDK
Simple SDK to interact with BeamNG.drive OutGauge data.

## Development
### Performance
`go test -bench=BenchmarkUpdate -benchmem -memprofile=mem.pprof`
replace the function to be tested

Use `go tool pprof` to analyze the results

#### Results
```bash
# Previous footprint
goos: linux
goarch: amd64
pkg: github.com/ESilva15/gobngsdk
cpu: AMD Ryzen 7 5800X3D 8-Core Processor
BenchmarkReadData-16    	 362514	     3188 ns/op	      4 B/op	      1 allocs/op
PASS
ok  	github.com/ESilva15/gobngsdk	1.194s

# New footprint
goos: linux
goarch: amd64
pkg: github.com/ESilva15/gobngsdk
cpu: AMD Ryzen 7 5800X3D 8-Core Processor
BenchmarkUpdate-16    	 159170	     7287 ns/op	      4 B/op	      1 allocs/op
PASS
ok  	github.com/ESilva15/gobngsdk	1.241s
# Pretty good enough. I can finally go be productive instead of "procrastinating" here
```
