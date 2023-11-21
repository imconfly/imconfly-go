package queue

import (
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"sync"
)

type Task struct {
	Request   *Request
	Origin    *Origin
	Transform *Transform
}

type Queue struct {
	mu    sync.Mutex
	tsMap map[os_tools.FileRelativePath][]chan error
	queue chan *Task
}

func NewQueue() *Queue {
	return &Queue{
		mu:    sync.Mutex{},
		tsMap: make(map[os_tools.FileRelativePath][]chan error),
		queue: make(chan *Task),
	}
}

// Close queue channel
// From now Get() returns nil and Add() throws panic on new tasks
func (q *Queue) Close() {
	close(q.queue)
}

func (q *Queue) Add(task *Task) chan error {
	q.mu.Lock()
	defer q.mu.Unlock()

	ch := make(chan error, 1)

	if alreadyAdded, found := q.tsMap[task.Request.Key]; found {
		alreadyAdded = append(alreadyAdded, ch)
	} else {
		q.tsMap[task.Request.Key] = []chan error{ch}
		q.queue <- task
	}

	return ch
}

func (q *Queue) Get() *Task {
	return <-q.queue
}

func (q *Queue) TaskDone(key os_tools.FileRelativePath, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if taskSubscribers, ok := q.tsMap[key]; ok {
		for _, ch := range taskSubscribers {
			ch <- err
		}
	} else {
		panic(fmt.Sprintf("Task %q not found in tasks map!", key))
	}

	delete(q.tsMap, key) // @todo: all read?
}
