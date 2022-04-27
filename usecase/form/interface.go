package form

import "Atlan_Collect_Challenge/entity"

type Reader interface {
	Extract(formId int8) (*entity.Form, error)
}

type Writer interface {
	Add(form *entity.Form) (bool, int, error)
}

type Deleter interface {
	Delete(formId int8) (bool, error)
}

type Repository interface {
	Reader
	Writer
	Deleter
}

type Usecase interface {
	GetForm(formId int8) (*entity.Form, error)
	CreateForm(form *entity.Form) (bool, int, error)
	DeleteForm(formId int8) (bool, error)
}
