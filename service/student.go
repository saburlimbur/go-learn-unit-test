package service

import (
	"session-9/model"
	"session-9/repository"
	"session-9/utils"
)

// StudentService interface used by CLI handler
type StudentServiceInterface interface {
	GetAll() ([]model.Student, error)
	GetByID(id int) (*model.Student, error)
	Create(student model.Student) (model.Student, error)
	Update(id int, student model.Student) (model.Student, error)
	Delete(id int) error
}

type StudentService struct {
	repo repository.StudentRepositoryInterface
}

func NewStudentService(repo repository.StudentRepositoryInterface) *StudentService {
	return &StudentService{repo: repo}
}

// GetAll returns all students
func (s *StudentService) GetAll() ([]model.Student, error) {
	// return s.repo.GetAll()
	students, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}
	if students == nil {
		return []model.Student{}, nil // memastikan tdk return nil slice
	}

	return students, nil
}

// GetByID returns a student by id
func (s *StudentService) GetByID(id int) (*model.Student, error) {
	students, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, st := range students {
		if st.ID == id {
			copy := st
			return &copy, nil
		}
	}
	return nil, utils.ErrNotFound
}

// Create creates new student with auto-increment ID
func (s *StudentService) Create(input model.Student) (model.Student, error) {
	students, err := s.repo.GetAll()
	if err != nil {
		return model.Student{}, err
	}

	maxID := 0
	for _, st := range students {
		if st.ID > maxID {
			maxID = st.ID
		}
	}
	input.ID = maxID + 1

	students = append(students, input)

	if err := s.repo.SaveAll(students); err != nil {
		return model.Student{}, err
	}

	return input, nil
}

// Update updates existing student by id
func (s *StudentService) Update(id int, input model.Student) (model.Student, error) {
	students, err := s.repo.GetAll()
	if err != nil {
		return model.Student{}, err
	}

	updated := false
	for i, st := range students {
		if st.ID == id {
			input.ID = id
			students[i] = input
			updated = true
			break
		}
	}

	if !updated {
		return model.Student{}, utils.ErrNotFound
	}

	if err := s.repo.SaveAll(students); err != nil {
		return model.Student{}, err
	}

	return input, nil
}

// Delete deletes student by id
func (s *StudentService) Delete(id int) error {
	students, err := s.repo.GetAll()
	if err != nil {
		return err
	}

	newList := make([]model.Student, 0, len(students))
	found := false
	for _, st := range students {
		if st.ID == id {
			found = true
			continue
		}
		newList = append(newList, st)
	}

	if !found {
		return utils.ErrNotFound
	}

	return s.repo.SaveAll(newList)
}
