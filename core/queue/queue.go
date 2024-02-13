package queue

import (
	"fmt"
	"github.com/imconfly/imconfly_go/core/origin"
	"github.com/imconfly/imconfly_go/core/request"
	"github.com/imconfly/imconfly_go/core/transform"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Task struct {
	Request   *request.Request
	Origin    *origin.Origin
	Transform *transform.Transform
}

type Queue struct {
	mu    *sync.Mutex
	tsMap map[os_tools.FileRelativePath][]chan error
	queue chan *Task
}

func NewQueue() *Queue {
	return &Queue{
		mu:    new(sync.Mutex),
		tsMap: make(map[os_tools.FileRelativePath][]chan error),
		queue: make(chan *Task),
	}
}

// Close queue channel
// From now Get() returns nil and Add() throws panic on new tasks
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	close(q.queue)
}

func (q *Queue) Add(task *Task) chan error {
	logName := fmt.Sprintf("Queue.Add(%s)", task.Request.Key)

	q.mu.Lock()
	defer q.mu.Unlock()

	ch := make(chan error, 1)

	if oldSlice, ok := q.tsMap[task.Request.Key]; ok {
		// just add subscriber
		oldSlice = append(oldSlice, ch)
		log.Debugf("%s: task already exists, just add subscriber.", logName)
	} else {
		// add new task and add subscriber
		q.tsMap[task.Request.Key] = []chan error{ch}
		q.queue <- task
		log.Debugf("%s: add new task.", logName)
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
