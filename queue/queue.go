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

type taskSubscribers struct {
	Task        *Task
	Subscribers []chan error
}

type Queue struct {
	mu    sync.Mutex
	tsMap map[os_tools.FileRelativePath]*taskSubscribers
	queue chan *Task
}

func NewQueue() *Queue {
	return &Queue{
		mu:    sync.Mutex{},
		tsMap: make(map[os_tools.FileRelativePath]*taskSubscribers),
		queue: make(chan *Task),
	}
}

// Close queue channel
// From now Get() returns nil and Add() throws panic on new tasks
func (q *Queue) Close() {
	close(q.queue)
}

func (q *Queue) Add(ta *Task) chan error {
	q.mu.Lock()
	defer q.mu.Unlock()

	ch := make(chan error, 1)

	if alreadyAdded, found := q.tsMap[ta.Request.Key]; found {
		alreadyAdded.Subscribers = append(alreadyAdded.Subscribers, ch)
	} else {
		ts := &taskSubscribers{
			Task:        ta,
			Subscribers: []chan error{ch},
		}
		q.tsMap[ta.Request.Key] = ts
		q.queue <- ta
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
		for _, ch := range taskSubscribers.Subscribers {
			ch <- err
		}
	} else {
		panic(fmt.Sprintf("Task `%s` not found in transform transform map!", key))
	}

	delete(q.tsMap, key) // @todo: all read?
}
