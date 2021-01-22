package tools

import (
	"sync"
	"time"
)

//Timer 定时器类
type Timer struct {
	//many timer
	timers map[int]*time.Timer

	//seed
	seed int

	//signal to exit
	exit map[int]chan struct{}

	//control go
	wg *sync.WaitGroup

	//lock for timers&exit map
	lock *sync.Mutex
}

//NewTimer 创建一个定时器
func NewTimer() *Timer {
	p := &Timer{}
	p.init()
	return p
}

//Init 初始化
func (t *Timer) init() {
	t.exit = make(map[int]chan struct{})
	t.wg = &sync.WaitGroup{}
	t.lock = &sync.Mutex{}
	t.timers = make(map[int]*time.Timer)
	t.seed = 0
}

//SetTimer 新建一个定时器,返回定时器id
func (t *Timer) SetTimer(d time.Duration, f func() bool) int {

	//创建一个定时器
	timer := time.NewTimer(d)
	//创建一个退出信号
	exit := make(chan struct{})

	t.seed++

	go func(timerid int) {
		t.wg.Add(1)
		defer t.wg.Done()

		for {
			select {
			case <-timer.C:
				if !f() {
					t.lock.Lock()
					defer t.lock.Unlock()

					if timer, ok := t.timers[t.seed]; ok {
						timer.Stop()
						delete(t.timers, timerid)
					}
					delete(t.exit, timerid)
					close(exit)
					goto End
				}
				timer.Reset(d)
			case <-exit:
				close(exit)
				goto End
			}
		}
	End:
		//log.Println("timer: ", timerid, "exit!")
	}(t.seed)

	t.lock.Lock()
	t.timers[t.seed] = timer
	t.exit[t.seed] = exit
	t.lock.Unlock()

	return t.seed
}

//Close 关闭定时器,0:全部关闭
func (t *Timer) Close(timerid int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if 0 == timerid {
		for _, v := range t.exit {
			v <- struct{}{}
		}

		for _, v := range t.timers {
			v.Stop()
		}

		t.exit = nil
		t.timers = nil

		t.wg.Wait()

		return
	}

	if exit, ok := t.exit[timerid]; ok {
		exit <- struct{}{}
		if timer, ok := t.timers[t.seed]; ok {
			timer.Stop()
			delete(t.timers, timerid)
		}
		delete(t.exit, timerid)
	}
}
