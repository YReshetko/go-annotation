.PHONY: bench
bench:
	go test -run=. -bench=. -benchtime=1s -count 1  -benchmem -cpuprofile=cpu.out -memprofile=mem.out -trace=trace.out bench_test.go | tee bench.txt

# go tool pprof -http :8080 cpu.out
# go tool pprof -http :8081 mem.out
# go tool trace trace.out
#
# go tool pprof $FILENAME.test cpu.out
# # (pprof) list <func name>
#
# # go get -u golang.org/x/perf/cmd/benchstat
# benchstat bench.txt
# rm cpu.out mem.out trace.out bench.txt *.test