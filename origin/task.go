package origin

type Key string

type Task struct {
	Remote string
	Key    Key // also local relative path
}
