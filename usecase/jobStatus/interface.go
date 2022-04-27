package jobStatus

import "Atlan_Collect_Challenge/entity"

type Reader interface {
	Extract(jobId int) (*entity.Job, error)
}

type Writer interface {
	Add(pluginCode int) (int, error)
	Update(jobId, statusCode int, status string) error
}

type Repository interface {
	Reader
	Writer
}

type Usecase interface {
	AddJob(pluginCode int) (int, error)
	UpdateJob(jobId, statusCode int, status string) error
	GetJob(jobId int) (*entity.Job, error)
}
