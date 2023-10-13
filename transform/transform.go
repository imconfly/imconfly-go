package transform

type Origin struct {
	Remote string
	Local  string
}

type Transform struct {
	Transform string
	Local     string
}

// Task - binding for request, transforms_conf and execution info
type Task struct {
	TaskRequest *TaskRequest
	Origin      *Origin
	Transform   *Transform
	Done        chan bool
}

func NewTask(tRec *TaskRequest, o *Origin, tr *Transform) *Task {
	return &Task{
		TaskRequest: tRec,
		Origin:      o,
		Transform:   tr,
		Done:        make(chan bool, 1),
	}
}
