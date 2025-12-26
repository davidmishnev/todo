package tasks

type TaskStatus int

const (
	Assigned TaskStatus = iota
	InProgress
	Completed
	Dropped
)

type Task struct {
	TaskID      int
	Header      string
	Description string
	Status      TaskStatus
}
