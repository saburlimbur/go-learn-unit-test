package service

import (
	"errors"
	"session-9/model"
	"session-9/repository"
	"session-9/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func newTestService(initial []model.Student) (*StudentService, *repository.MockStudentRepository) {
// 	repo := &repository.MockStudentRepository{Students: initial}
// 	svc := NewStudentService(repo)
// 	return svc, repo
// }

// func TestStudentService_Create(t *testing.T) {
// 	svc, repo := newTestService([]model.Student{})

// 	created, err := svc.Create(model.Student{
// 		Name: "Budi",
// 		Age:  20,
// 	})
// 	if err != nil {
// 		t.Fatalf("Create returned error: %v", err)
// 	}

// 	if created.ID != 1 {
// 		t.Errorf("expected ID 1, got %d", created.ID)
// 	}
// 	if created.Name != "Budi" {
// 		t.Errorf("expected Name Budi, got %s", created.Name)
// 	}

// 	if len(repo.Students) != 1 {
// 		t.Fatalf("expected repo to have 1 student, got %d", len(repo.Students))
// 	}
// }

// func TestStudentService_GetByID_Found(t *testing.T) {
// 	initial := []model.Student{
// 		{ID: 1, Name: "Andi", Age: 21},
// 		{ID: 2, Name: "Siti", Age: 22},
// 	}
// 	svc, _ := newTestService(initial)

// 	st, err := svc.GetByID(2)
// 	if err != nil {
// 		t.Fatalf("GetByID returned error: %v", err)
// 	}

// 	if st.Name != "Siti" {
// 		t.Errorf("expected Name Siti, got %s", st.Name)
// 	}
// }

// func TestStudentService_GetByID_NotFound(t *testing.T) {
// 	initial := []model.Student{
// 		{ID: 1, Name: "Andi", Age: 21},
// 		{ID: 2, Name: "Siti", Age: 22},
// 	}
// 	svc, _ := newTestService(initial)

// 	_, err := svc.GetByID(999)
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}

// 	if err != utils.ErrNotFound {
// 		t.Fatalf("expected ErrNotFound, got %v", err)
// 	}
// }

// func TestStudentService_GetByID_fileError(t *testing.T) {
// 	svc, _ := newTestService([]model.Student{})

// 	_, err := svc.GetByID(1)
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}

// 	if err != utils.ErrFile {
// 		t.Fatalf("expected error file, got %v", err)
// 	}
// }

func newTestService() (*StudentService, *repository.MockStudentRepository) {
	mokeRepo := new(repository.MockStudentRepository)
	service := NewStudentService(mokeRepo)
	return service, mokeRepo
}

func TestStudentService_Create_Success_File_err(t *testing.T) {
	svc, repo := newTestService()

	// existing := []model.Student{
	// 	{ID: 3, Name: "Andi", Age: 21},
	// 	{ID: 5, Name: "Siti", Age: 22},
	// }

	repo.On("GetAll").Return([]model.Student{}, utils.ErrFile).Once()
	// repo.On("SaveAll").Return(nil).Once()

	input := model.Student{
		Name: "Budi",
		Age:  20,
	}

	_, err := svc.Create(input)

	assert.Error(t, err)
	// assert.Equal(t, 6, created.ID)
	// assert.Equal(t, "Budi", created.Name)
	assert.Equal(t, utils.ErrFile, err)

	repo.AssertExpectations(t)
}

func TestStudentService_Create_Error_SaveAll(t *testing.T) {
	svc, repo := newTestService()

	existing := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
	}

	repo.
		On("GetAll").
		Return(existing, nil).
		Once()

	appendData := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Budi", Age: 20},
	}

	repo.
		On("SaveAll", appendData).
		Return(errors.New("save error")).
		Once()

	input := model.Student{
		Name: "Budi",
		Age:  20,
	}

	created, err := svc.Create(input)

	assert.Error(t, err)
	assert.EqualError(t, err, "save error")
	assert.Equal(t, model.Student{}, created)

	repo.AssertExpectations(t)
}

func TestStudentService_GetByID_Found(t *testing.T) {
	initial := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Siti", Age: 22},
	}
	svc, repo := newTestService()
	repo.On("GetAll").Return(initial, nil).Once()

	st, err := svc.GetByID(2)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if st.Name != "Siti" {
		t.Errorf("expected Name Siti, got %s", st.Name)
	}
}

func TestStudentService_GetByID_NotFound(t *testing.T) {
	initial := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Siti", Age: 22},
	}
	svc, repo := newTestService()
	repo.On("GetAll").Return(initial, nil).Once()

	_, err := svc.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)

	repo.AssertExpectations(t)
}

// func TestStudentService_GetByID_NotFound(t *testing.T) {
// 	initial := []model.Student{
// 		{ID: 1, Name: "Andi", Age: 21},
// 		{ID: 2, Name: "Siti", Age: 22},
// 	}
// 	svc, repo := newTestService()
// 	repo.On("GetAll").Return(initial, utils.ErrFile).Once()

// 	_, err := svc.GetByID(999)
// 	if err == nil {
// 		t.Fatalf("expected error, got nil")
// 	}

// 	if err != utils.ErrNotFound {
// 		t.Fatalf("expected ErrNotFound, got %v", err)
// 	}
// }

func TestStudentService_GetByID_FileError(t *testing.T) {
	svc, repo := newTestService()
	repo.On("GetAll").Return([]model.Student{}, utils.ErrFile).Once()

	_, err := svc.GetByID(1)

	assert.Error(t, err)
	assert.Equal(t, utils.ErrFile, err)
}

func TestStudentService_Update(t *testing.T) {
	svc, repo := newTestService()

	existing := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Budi", Age: 20},
	}

	input := model.Student{
		Name: "Budi Dermawan",
		Age:  20,
	}

	t.Run("student file error", func(t *testing.T) {
		repo.
			On("GetAll").
			Return([]model.Student{}, utils.ErrFile).
			Once()

		created, err := svc.Update(2, input)

		assert.Error(t, err)
		assert.EqualError(t, err, "file error")
		assert.Equal(t, model.Student{}, created)

		repo.AssertExpectations(t)
	})

	t.Run("student not found", func(t *testing.T) {
		repo.
			On("GetAll").
			Return(existing, nil).
			Once()

		created, err := svc.Update(3, input)

		assert.Error(t, err)
		assert.EqualError(t, err, "student not found")
		assert.Equal(t, model.Student{}, created)

		repo.AssertExpectations(t)
	})

}
