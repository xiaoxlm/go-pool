# gopool
a simple goroutine pool
# Feature
1. support sync and async process
2. when idle worker(gorouine) is timeout, it will quit
3. add worker to max numbers auto-dynamically 
# Example
usage:
```go
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
```
more usage examples in test file.