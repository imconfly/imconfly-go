package origin

import (
	"fmt"
	"sync"
)

type taskSubscribers struct {
	Task        *Task
	Subscribers []chan error
}

type Queue struct {
	mu                 sync.Mutex
	taskSubscribersMap map[Key]*taskSubscribers
	queueChan          chan *Task
}

func (q *Queue) Add(ta *Task, subscriber chan error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if alreadyAdded, found := q.taskSubscribersMap[ta.Key]; found {
		alreadyAdded.Subscribers = append(alreadyAdded.Subscribers, subscriber)
	} else {
		taskSubscribers := &taskSubscribers{
			Task:        ta,
			Subscribers: []chan error{subscriber},
		}
		q.taskSubscribersMap[ta.Key] = taskSubscribers
		q.queueChan <- ta
	}
}

func (q *Queue) Get() *Task {
	return <-q.queueChan
}

func (q *Queue) Done(key Key, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if taskSubscribers, ok := q.taskSubscribersMap[key]; ok {
		for _, ch := range taskSubscribers.Subscribers {
			ch <- err
		}
	} else {
		panic(fmt.Sprintf("Task `%s` not found in origin transform map!", key))
	}

	delete(q.taskSubscribersMap, key)
}
