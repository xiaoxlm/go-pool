package gopool

import (
	"sync"
	"testing"
	"time"
)

const (
	RunTimes           = 1000
	BenchParam         = 10
	BenchAntsSize      = 200000
	DefaultExpiredTime = 10 * time.Second
	JobBuffer          = 30
)

func demoFunc() {
	//time.Sleep(1 * time.Millisecond)
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}


func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkPool(b *testing.B) {
	var wg sync.WaitGroup

	gopool, _ := NewPool(Options{
		MaxWorkerNums:   BenchAntsSize,
		MinWorkerNums:   1,
		JobBuffer:   JobBuffer,
		IdleTimeout: DefaultExpiredTime,
	})

	defer gopool.Close()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			_ = gopool.ProcessAsync(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}