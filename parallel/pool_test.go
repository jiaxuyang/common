package parallel

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestPool_Success(t *testing.T) {
	add := func(a int, b int) (result int, err error) {
		time.Sleep(1 * time.Second)
		return a + b, nil
	}
	// total job count
	N := 100
	as := make([]int, N)
	bs := make([]int, N)
	results := make([]int, N)
	errs := make([]error, N)
	fns := make([]func(), 0, N)
	for i := 0; i < N; i++ {
		as[i] = i + 1
		bs[i] = N - i
	}
	fnC := make(chan func(), N)
	wg := new(sync.WaitGroup)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i2 int) {
			fnC <- func() {
				results[i2], errs[i2] = add(as[i2], bs[i2])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	// must close?
	close(fnC)
	for fn := range fnC {
		fns = append(fns, fn)
	}
	log.Println("fns", len(fns))
	pool := NewPool(&Options{MaxPoolSize: 5})
	for _, fn := range fns {
		pool.Go(fn)
	}
	pool.Wait()
	log.Println("finish")
	for i, r := range results {
		if r != N+1 {
			t.Errorf("result error. index=%d result=%v", i, r)
		}
	}
	for i, err := range errs {
		if err != nil {
			t.Errorf("errs error. index=%d err=%v", i, err)
		}
	}
}

func TestPool_Panic(t *testing.T) {
	add := func(a int, b int) (result int, err error) {
		time.Sleep(1 * time.Millisecond)
		panic(a + b)
		// return a + b, nil
	}
	N := 10
	as := make([]int, N)
	bs := make([]int, N)
	results := make([]int, N)
	errs := make([]error, N)
	fns := make([]func(), 0, N)
	for i := 0; i < N; i++ {
		as[i] = i + 1
		bs[i] = N - i
	}
	fnC := make(chan func(), N)
	wg := new(sync.WaitGroup)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i2 int) {
			fnC <- func() {
				// user handle panic
				defer func() {
					if r := recover(); r != nil {
						log.Printf("recover %d", r)
					}
				}()
				results[i2], errs[i2] = add(as[i2], bs[i2])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	// must close?
	close(fnC)
	for fn := range fnC {
		fns = append(fns, fn)
	}
	pool := NewPool(&Options{MaxPoolSize: 5})
	for _, fn := range fns {
		pool.Go(fn)
	}
	pool.Wait()
	log.Println("finish")
	for i, r := range results {
		if r != 0 {
			t.Errorf("result error. index=%d result=%v", i, r)
		}
	}
	for i, err := range errs {
		if err != nil {
			t.Errorf("errs error. index=%d err=%v", i, err)
		}
	}
}
