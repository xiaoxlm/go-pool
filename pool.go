package gopool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_MAX_WORKER   = 100
	DEFAULT_IDLE_TIMEOUT = 60 * time.Second

	BUCKETSIZE = 5
)

type jobFunc func()

type task struct {
	result chan bool
	call   jobFunc
}

type Options struct {
	MaxWorkerNums int32
	MinWorkerNums int32
	IdleTimeout   time.Duration
	JobBuffer     int
}

type Pool struct {
	sync.Mutex

	maxWorkerNums int32
	minWorkerNums int32
	curWorkerNums int32
	idleTimeout   time.Duration

	jobChan   chan task
	jobBuffer int

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewPool(option Options) (*Pool, error) {
	pool := Pool{
		maxWorkerNums: option.MaxWorkerNums,
		minWorkerNums: option.MinWorkerNums,
		idleTimeout: option.IdleTimeout,
		jobBuffer: option.JobBuffer,
	}

	if pool.maxWorkerNums < pool.minWorkerNums {
		return nil, errors.New("minWokerNums must less than maxWorkerNums")
	}

	if pool.maxWorkerNums < 1 {
		pool.maxWorkerNums = DEFAULT_MAX_WORKER
	}

	if pool.minWorkerNums < 1 {
		pool.minWorkerNums = (pool.maxWorkerNums / BUCKETSIZE) + 1
	}

	if pool.idleTimeout < time.Second {
		pool.idleTimeout = DEFAULT_IDLE_TIMEOUT
	}

	pool.jobChan = make(chan task, pool.jobBuffer)
	pool.ctx, pool.ctxCancel = context.WithCancel(context.Background())
	pool.addWorker(int(pool.maxWorkerNums))

	return &pool, nil
}

func (p *Pool) addWorker(nums int) {
	for n := 0; n < nums; n++ {
		if p.curWorkerNums >= p.maxWorkerNums {
			return
		}
		atomic.AddInt32(&(p.curWorkerNums), 1)
		go p.runGoroutine()
	}
}

func (p *Pool) deleteWorker() error {
	if p.curWorkerNums <= p.minWorkerNums {
		return errors.New("can't less minWorkerNums")
	}
	atomic.AddInt32(&(p.curWorkerNums), -1)
	return nil
}

// exec goroutine
func (p *Pool) runGoroutine() {
	timer := time.NewTimer(p.idleTimeout)
	for {
		select {
		case task := <-p.jobChan:
			task.call()
			if task.result != nil {
				task.result <- true
			}
			timer.Reset(p.idleTimeout)
		case <-timer.C:
			if err := p.deleteWorker(); err != nil {
				timer.Reset(p.idleTimeout)
				continue
			}
			return
		case <-p.ctx.Done():
			if p.curWorkerNums > 0 {
				atomic.AddInt32(&(p.curWorkerNums), -1)
			}
			return
		}
	}
}

func (p *Pool) completeWorker() {
	var addWorkerNums int
	addWorkerNums = int(p.maxWorkerNums) / BUCKETSIZE

	p.addWorker(addWorkerNums)
}

func (p *Pool) ProcessAsync(f jobFunc) error {
	t := task{
		call:   f,
		result: nil,
	}

	p.completeWorker()
	p.jobChan <- t
	return nil
}

func (p *Pool) ProcessSync(f jobFunc) error {
	t := task{
		call:   f,
		result: make(chan bool, 0),
	}

	p.completeWorker()
	p.jobChan <- t
	<-t.result
	return nil
}

func (p *Pool) GetCurWorkerNums() int32 {
	return p.curWorkerNums
}

func (p *Pool) Close() {
	p.ctxCancel()
	return
}
