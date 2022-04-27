package entity

type Form struct {
	FormName string     `json:"form_name" binding:"required"`
	FormId   int8       `json:"form_id"`
	OwnerId  int        `json:"owner_id" binding:"required"`
	Question []Question `json:"question" binding:"required"`
}

func NewForm(name string, ownerId int, questions []Question) Form {

	return Form{
		FormName: name,
		OwnerId:  ownerId,
		Question: questions,
	}
}
