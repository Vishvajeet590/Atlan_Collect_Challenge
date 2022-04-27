package response

import "Atlan_Collect_Challenge/entity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetResponses(formId int8) (*[]entity.Form, error) {
	return nil, nil
}

func (s *Service) AddResponse(response *entity.Response, formId, userId int8) (bool, error) {
	res, err := s.repo.Add(response, formId, userId)
	if err != nil {
		return false, err
	}
	return res, nil

}
