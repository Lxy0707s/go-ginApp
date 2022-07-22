package queue

import (
	"fmt"
	"runtime"
	"time"
)

func addFunc(q *CASQueue, prefix int, l int) {
	for i := 0; i < l; i++ {
		v := fmt.Sprintf("%d---%d", prefix, i)
		err := q.Put(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v ", v)
	}
	err := q.Put("0---5")
	if err != nil {
		fmt.Println(err)
	}
	err = q.Put("1---5")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("0---5 1---5")
	fmt.Println("end")
}

func getFunc(q *CASQueue, l int) {
	v, flag := q.Get()
	if !flag {
		fmt.Printf("get fail, the queue is empty\n")
	}
	fmt.Println(v)
}

// 测试并发时，查看是否做到了线程安全
func TestQueue() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	q := NewCASQueue(10000)
	l := 10
	go func() {
		for i := 0; i < 2; i++ {
			prefix := i
			fmt.Printf("truely prefix [%d]\n", prefix)
			addFunc(q, prefix, l)
		}
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("going")
	fmt.Println("len", q.Quantity())

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		/*case <-c.TaskResource.SuccessSign:
		c.startChan <- 1
		c.log.Info("任务更新，触发任务分配")*/
		case <-ticker.C:
			for {
				time.Sleep(1 * time.Second)
				fmt.Println("going")
				time.Sleep(1 * time.Second)
				fmt.Println(q.Quantity())
				if q.Quantity() != 0 {
					l = int(q.Quantity())
					getFunc(q, l)
				}
			}

		}
	}
}
