package form

import "Atlan_Collect_Challenge/entity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetForm(formId int8) (*entity.Form, error) {
	form, err := s.repo.Extract(formId)
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (s *Service) CreateForm(form *entity.Form) (bool, int, error) {
	res, formId, err := s.repo.Add(form)
	if err != nil || res == false {
		return false, -999, err
	}
	return res, formId, nil
}

func (s *Service) DeleteForm(formId int8) (bool, error) {
	return true, nil
}
