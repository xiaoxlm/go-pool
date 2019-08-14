package gopool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	wg := sync.WaitGroup{}
	gopool, err := NewPool(Options{
		MaxWorkerNums:   10,
		MinWorkerNums:   3,
		JobBuffer:   5,
		IdleTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	for index := 0; index < 1000; index++ {
		wg.Add(1)
		f := func(){
			time.Sleep(1 * time.Second)
			fmt.Println(time.Now(), "::666", "worker nums:", gopool.GetCurWorkerNums())
			wg.Done()
		}

		if err := gopool.ProcessAsync(f); err != nil {
			fmt.Println("err:", err)
		}
	}

	wg.Wait()
}

func TestGetCurWorkerNums(t *testing.T) {
	wg := sync.WaitGroup{}
	gopool, err := NewPool(Options{
		MaxWorkerNums:   10,
		MinWorkerNums:   3,
		JobBuffer:   1,
		IdleTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	for index := 0; index < 5; index++ {
		wg.Add(1)
		f := func(){
			time.Sleep(1 * time.Second)
			fmt.Println(time.Now(), "::666", "worker nums:", gopool.GetCurWorkerNums())
			wg.Done()
		}

		if err := gopool.ProcessAsync(f); err != nil {
			fmt.Println("err:", err)
		}
	}

	wg.Wait()

	for{
		fmt.Println(gopool.GetCurWorkerNums())
		time.Sleep(3 * time.Second)
	}
}

func TestCloseAllWorkers(t *testing.T) {
	wg := sync.WaitGroup{}
	gopool, err := NewPool(Options{
		MaxWorkerNums:   10,
		MinWorkerNums:   3,
		JobBuffer:   1,
		IdleTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	for index := 0; index < 5; index++ {
		wg.Add(1)
		f := func(){
			time.Sleep(1 * time.Second)
			fmt.Println(time.Now(), "::666", "worker nums:", gopool.GetCurWorkerNums())
			wg.Done()
		}

		if err := gopool.ProcessAsync(f); err != nil {
			fmt.Println("err:", err)
		}
	}

	wg.Wait()

	go func(){
		time.Sleep(20 * time.Second)
		gopool.Close()
	}()

	for{
		fmt.Println(gopool.GetCurWorkerNums())
		time.Sleep(3 * time.Second)
	}
}