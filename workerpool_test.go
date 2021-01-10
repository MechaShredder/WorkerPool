package workerpool

import (
	"sync"
	"testing"
)

const taskNum = 10

var result = make([]int, taskNum)
var taskList = make([]*CountTask, taskNum)
var pool = NewWorkerPool(taskNum)

func init() {
	for i := 0; i < taskNum; i++ {
		taskList[i] = NewCountTask(i)
	}
	pool.Start()
}

type CountTask struct {
	id    int
	count int
	quit  chan struct{}
}

func NewCountTask(id int) *CountTask {
	return &CountTask{id: id, quit: make(chan struct{})}
}

func (t *CountTask) Do() {
	t.count++
	result[t.id] = t.count
	t.quit <- struct{}{}
}

func (t *CountTask) wait() {
	<-t.quit
}

func TestPerformance(t *testing.T) {
	count := 2
	for i := 0; i < count; i++ {
		t.Logf("test count: %d\n", i)
		ConcurrencyTest(taskNum)
	}
	Check(count, t)
}

func ConcurrencyTest(size int) {
	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(id int) {
			pool.Assign(taskList[id])
			taskList[id].wait()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func Check(target int, t *testing.T) {
	for i := 0; i < taskNum; i++ {
		if result[i] != target {
			t.Errorf("result[%d] = %d, failed", i, result[i])
			return
		}
	}
}
