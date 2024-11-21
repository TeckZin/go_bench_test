package bench

import (
	"fmt"
	"runtime"
	"sort"
	"time"
)

// BenchmarkResult stores the results of a single benchmark run
type BenchmarkResult struct {
	Name        string
	TotalTime   time.Duration
	Iterations  int
	AverageTime time.Duration
	MemoryUsage uint64
	AllocCount  uint64
}

// Benchmark is the main benchmarking utility
type Benchmark struct {
	results []BenchmarkResult
}

// measurement represents a single iteration measurement
type measurement struct {
	duration    time.Duration
	memoryDelta uint64
	allocCount  uint64
}

// New creates a new Benchmark instance
func New() *Benchmark {
	return &Benchmark{
		results: make([]BenchmarkResult, 0),
	}
}

// Run executes a single benchmark
func (b *Benchmark) Run(name string, fn func(), iterations int) *Benchmark {
	fmt.Printf("Starting: %s\n", name)

	// Force GC before starting
	runtime.GC()

	// Warm up run
	fn()

	measurements := make([]measurement, iterations)

	for i := 0; i < iterations; i++ {
		// Get memory stats before
		var memStatsBefore runtime.MemStats
		runtime.ReadMemStats(&memStatsBefore)

		start := time.Now()
		fn()
		duration := time.Since(start)

		// Get memory stats after
		var memStatsAfter runtime.MemStats
		runtime.ReadMemStats(&memStatsAfter)

		measurements[i] = measurement{
			duration:    duration,
			memoryDelta: memStatsAfter.HeapAlloc - memStatsBefore.HeapAlloc,
			allocCount:  memStatsAfter.Mallocs - memStatsBefore.Mallocs,
		}
	}

	// Calculate averages
	var totalDuration time.Duration
	var totalMemory uint64
	var totalAllocs uint64

	for _, m := range measurements {
		totalDuration += m.duration
		totalMemory += m.memoryDelta
		totalAllocs += m.allocCount
	}

	avgDuration := totalDuration / time.Duration(iterations)
	avgMemory := totalMemory / uint64(iterations)
	avgAllocs := totalAllocs / uint64(iterations)

	b.results = append(b.results, BenchmarkResult{
		Name:        name,
		TotalTime:   totalDuration,
		Iterations:  iterations,
		AverageTime: avgDuration,
		MemoryUsage: avgMemory,
		AllocCount:  avgAllocs,
	})

	fmt.Printf("Done: %s\n", name)
	return b
}

// Compare runs multiple benchmarks
func (b *Benchmark) Compare(tests map[string]func(), iterations int) *Benchmark {
	for name, fn := range tests {
		b.Run(name, fn, iterations)
	}
	return b
}

// PrintResults displays the benchmark results in a formatted table
func (b *Benchmark) PrintResults() {
	fmt.Println("\nBenchmark Results:")
	fmt.Printf("%-20s %-15s %-15s %-15s %-15s %-15s\n",
		"Test Name", "Total Time", "Iterations", "Avg Time", "Memory (MB)", "Allocs")

	for _, result := range b.results {
		fmt.Printf("%-20s %-15s %-15d %-15s %-15.2f %-15d\n",
			result.Name,
			result.TotalTime,
			result.Iterations,
			result.AverageTime,
			float64(result.MemoryUsage)/(1024*1024),
			result.AllocCount)
	}

	// Find and report the fastest test
	sort.Slice(b.results, func(i, j int) bool {
		return b.results[i].AverageTime < b.results[j].AverageTime
	})

	fastest := b.results[0]
	fmt.Printf("\nFastest test: %s\n", fastest.Name)

	for _, result := range b.results[1:] {
		ratio := float64(result.AverageTime) / float64(fastest.AverageTime)
		fmt.Printf("%s is %.2fx slower than %s\n",
			result.Name, ratio, fastest.Name)
	}
}

// Clear resets the benchmark results
func (b *Benchmark) Clear() *Benchmark {
	b.results = b.results[:0]
	return b
}

// Example benchmark functions
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func SieveOfEratosthenes(n int) []int {
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// Example usage
func ExampleUsage() {
	bench := New()

	tests := map[string]func(){
		"fibonacci(45)": func() {
			Fibonacci(45)
		},
		"primes(1000000)": func() {
			SieveOfEratosthenes(1000000)
		},
	}

	bench.Compare(tests, 5).PrintResults()
}
