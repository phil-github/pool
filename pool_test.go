package pool

import (
	"fmt"
	"github.com/phil-github/pool"
	"net"
	"sync"
	"testing"
	"time"
)

func TestRemain(t *testing.T) {
	factory := func() (interface{}, error) {
		return net.Dial("tcp", "127.0.0.1:8000")
	}
	close := func(v interface{}) error {
		return v.(net.Conn).Close()
	}
	init := 2
	max := 5
	idle := 15

	poolConfig := &pool.PoolConfig{
		InitialCap:  init,
		MaxCap:      max,
		Factory:     factory,
		Close:       close,
		IdleTimeout: time.Duration(idle) * time.Second,
	}

	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println("err=", err)
	}

	size := p.Len()
	remain := p.Remain()

	if size != init || remain != max-init {
		t.Errorf("failed, init=%d, max=%d, len=%d, remain=%d", init, max, size, remain)
	}
	t.Logf("init=%d, max=%d, len=%d, remain=%d", init, max, size, remain)

	var connLst []interface{}
	try := 100
	for i := 1; i <= try; i++ {
		v, err := p.Get()

		size = p.Len()
		remain = p.Remain()

		if err == nil {
			connLst = append(connLst, v)
			if i+size+remain != max {
				t.Errorf("failed, err=%s, i=%d, init=%d, max=%d, len=%d, remain=%d", err, i, init, max, size, remain)
			}
		} else {
			if size != 0 || remain != 0 {
				t.Errorf("failed, i=%d, init=%d, max=%d, len=%d, remain=%d", i, init, max, size, remain)
			}
		}
	}

	for i, c := range connLst {
		p.Put(c)

		size = p.Len()
		remain = p.Remain()

		if i+1 != size || remain != 0 {
			t.Errorf("failed, i=%d, init=%d, max=%d, len=%d, remain=%d", i, init, max, size, remain)
		}
	}

	time.Sleep(time.Duration(idle+1) * time.Second)

	var connLst2 []interface{}
	for i := 1; i <= try; i++ {
		v, err := p.Get()

		size = p.Len()
		remain = p.Remain()

		if size != 0 {
			t.Errorf("failed, i=%d, init=%d, max=%d, len=%d, remain=%d", i, init, max, size, remain)
		}

		if err == nil {
			connLst2 = append(connLst2, v)
			if i+size+remain != max {
				t.Errorf("failed, err=%s, i=%d, init=%d, max=%d, len=%d, remain=%d", err, i, init, max, size, remain)
			}
		} else {
			if size != 0 || remain != 0 {
				t.Errorf("failed, i=%d, init=%d, max=%d, len=%d, remain=%d", i, init, max, size, remain)
			}
		}
	}

	for _, c := range connLst2 {
		p.Put(c)
	}

	p.Release()
}

func TestLimitPool(t *testing.T) {
	factory := func() (interface{}, error) {
		return net.Dial("tcp", "127.0.0.1:8000")
	}
	close := func(v interface{}) error {
		return v.(net.Conn).Close()
	}
	init := 1
	max := 200
	idle := 15

	poolConfig := &pool.PoolConfig{
		InitialCap:  init,
		MaxCap:      max,
		Factory:     factory,
		Close:       close,
		IdleTimeout: time.Duration(idle) * time.Second,
	}

	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		t.Errorf("err=%s", err)
	}

	size := p.Len()
	remain := p.Remain()

	if size != init || remain != max-init {
		t.Errorf("failed, init=%d, max=%d, len=%d, remain=%d", init, max, size, remain)
	}
	t.Logf("init=%d, max=%d, len=%d, remain=%d", init, max, size, remain)

	// consume
	wg := sync.WaitGroup{}
	wg.Add(2000)
	for i := 0; i < 2000; i++ {
		go func() {
			start := time.Now().Unix()

			for time.Now().Unix()-start < 30 {
				v, err := p.Get()
				time.Sleep(3 * time.Millisecond)
				if err == nil {
					p.Put(v)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
	size = p.Len()
	remain = p.Remain()

	if max != size+remain {
		t.Errorf("failed, max=%d, len=%d, remain=%d", max, size, remain)
	}
}
