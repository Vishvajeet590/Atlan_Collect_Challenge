package response

import "Atlan_Collect_Challenge/entity"

type Reader interface {
	Extract(formId int8) (*[]entity.Response, error)
}

type Writer interface {
	Add(form *entity.Response, formId, userId int8) (bool, error)
}

type Repository interface {
	Reader
	Writer
}

type Usecase interface {
	GetResponses(formId int8) (*[]entity.Response, error)
	AddResponse(response *entity.Response, formId, userId int8) (bool, error)
}
