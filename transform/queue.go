package transform

import (
	"fmt"
	"github.com/imconfly/imconfly_go/lib/os_tools"
	"sync"
)

type taskSubscribers struct {
	Task        *Task
	Subscribers []chan error
}

type Queue struct {
	mu        sync.Mutex
	tsMap     map[os_tools.FileRelativePath]*taskSubscribers
	queueChan chan *Task
}

func (q *Queue) Add(ta *Task, subscriberCh chan error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if alreadyAdded, found := q.tsMap[ta.File.Key]; found {
		alreadyAdded.Subscribers = append(alreadyAdded.Subscribers, subscriberCh)
	} else {
		ts := &taskSubscribers{
			Task:        ta,
			Subscribers: []chan error{subscriberCh},
		}
		q.tsMap[ta.File.Key] = ts
		q.queueChan <- ta
	}
}

func (q *Queue) Get() *Task {
	return <-q.queueChan
}

func (q *Queue) Done(key os_tools.FileRelativePath, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if taskSubscribers, ok := q.tsMap[key]; ok {
		for _, ch := range taskSubscribers.Subscribers {
			ch <- err
		}
	} else {
		panic(fmt.Sprintf("Task `%s` not found in transform transform map!", key))
	}

	delete(q.tsMap, key)
}
