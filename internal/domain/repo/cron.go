package repo

type Job struct {
	Name 	string
	Spec 	string
	Cmd 	func()
}

type CronRepository interface {
	JobList() map[string]string
	NewLocalJob(job *Job) error
	RunJob()
	RemoveJob(name string, id int) error
}
