package transform

type Origin struct {
	Source string
	Access bool
}

// Task - binding for request, transforms_conf and execution info
type Task struct {
	TaskRequest *TaskRequest
	Origin      *Origin
	Transform   *string
	Done        chan bool
}

func NewTask(tRec *TaskRequest, o *Origin, tr *string) *Task {
	return &Task{
		TaskRequest: tRec,
		Origin:      o,
		Transform:   tr,
		Done:        make(chan bool, 1),
	}
}
