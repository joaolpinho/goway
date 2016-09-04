package worker


type JobHandlerFunc func(job *Job)(bool)

type JobMap map[string]JobHandlerFunc

type ITaskWorker interface {
	Run(job *Job) bool
	AddJob(action string, handler JobHandlerFunc)
}

func (task *TaskWorker) AddJob(action string, handler JobHandlerFunc) {
	task.JobMap[action] = handler
}

func NewTaskWorker() *TaskWorker {
	return &TaskWorker{JobMap{}}
}

type TaskWorker struct{
	JobMap JobMap
}

func (task *TaskWorker) Run(job *Job) bool{
	fn := task.JobMap[job.Name]
	if fn == nil {
		return false
	}

	fn(job)

	return true
}