package queue

import (
	"fmt"
	"github.com/imconfly/imconfly_go/task"
	"sync"
)

type taskSubscribers struct {
	Task        *task.Task
	Subscribers []chan error
}

type Queue struct {
	mu        sync.Mutex
	tsMap     map[task.Key]*taskSubscribers
	queueChan chan *task.Task
}

func (q *Queue) Add(ta *task.Task, subscriberCh chan error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if alreadyAdded, found := q.tsMap[ta.Request.Key]; found {
		alreadyAdded.Subscribers = append(alreadyAdded.Subscribers, subscriberCh)
	} else {
		ts := &taskSubscribers{
			Task:        ta,
			Subscribers: []chan error{subscriberCh},
		}
		q.tsMap[ta.Request.Key] = ts
		q.queueChan <- ta
	}
}

func (q *Queue) Get() *task.Task {
	return <-q.queueChan
}

func (q *Queue) Done(key task.Key, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Уведомить всех подписчиков о выполнении
	{
		taskListeners, ok := q.tsMap[key]
		if !ok {
			panic(fmt.Sprintf("Task `%s` not found in queue map!", key))
		}
		for _, ch := range taskListeners.Subscribers {
			ch <- err
		}
	}

	delete(q.tsMap, key)
}
