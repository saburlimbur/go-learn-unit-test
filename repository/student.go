package repository

import (
	"session-9/model"
	"session-9/utils"
)

type StudentRepositoryInterface interface {
	GetAll() ([]model.Student, error)
	SaveAll([]model.Student) error
}

type StudentRepository struct {
	FilePath string
}

func NewStudentRepository(path string) *StudentRepository {
	return &StudentRepository{FilePath: path}
}

func (r *StudentRepository) GetAll() ([]model.Student, error) {
	var students []model.Student
	if err := utils.ReadJSON(r.FilePath, &students); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) SaveAll(students []model.Student) error {
	return utils.WriteJSON(r.FilePath, students)
}
