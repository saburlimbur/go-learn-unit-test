package service

import (
	"errors"
	"session-9/model"
	"session-9/repository"
	"session-9/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestService() (*StudentService, *repository.MockStudentRepository) {
	mokeRepo := new(repository.MockStudentRepository)
	service := NewStudentService(mokeRepo)
	return service, mokeRepo
}

// ======= BY ID FOUND
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
	// once hanya sekali panggil

	st, err := svc.GetByID(2)

	assert.NoError(t, err) // tidak mengembalikan error
	assert.NotNil(t, st)   // objek yang dikembalikan tidak nil

	// menyesuaikan isi data student seperti Siti
	assert.Equal(t, 2, st.ID)
	assert.Equal(t, "Siti", st.Name)
	assert.Equal(t, 22, st.Age)

	repo.AssertExpectations(t)
}

// ======= BY ID NOT FOUND
func TestStudentService_GetByID_NotFound(t *testing.T) {
	initial := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Siti", Age: 22},
	}
	svc, repo := newTestService()

	repo.On("GetAll").Return(initial, nil).Once()

	repo.On("GetAll").Return(initial, nil).Once()

	st, err := svc.GetByID(999)

	assert.Nil(t, st)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrNotFound, err)

	repo.AssertExpectations(t)
}

// ======= BY ID FILE ERROR
func TestStudentService_GetByID_fileError(t *testing.T) {
	svc, repo := newTestService() // instance service dan mock

	// ketika gagal return maka dia mengembalikan array dan error
	repo.On("GetAll").Return([]model.Student{}, utils.ErrFile).Once()

	_, err := svc.GetByID(1)

	assert.Error(t, err)
	assert.Equal(t, utils.ErrFile, err)

	repo.AssertExpectations(t)
}

// sub test
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

// table test
func TestStudentService_Delete(t *testing.T) {
	svc, repo := newTestService()

	tests := []struct {
		name         string
		paramID      int
		existingData []model.Student
		errGet       error
		messageError string
	}{
		{
			name:         "file error",
			paramID:      1,
			existingData: []model.Student{},
			errGet:       utils.ErrFile,
			messageError: "file error",
		},
		{
			name:    "not found",
			paramID: 3,
			existingData: []model.Student{
				{ID: 1, Name: "Andi", Age: 21},
				{ID: 2, Name: "Budi", Age: 20},
			},
			errGet:       utils.ErrNotFound,
			messageError: "student not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.
				On("GetAll").
				Return(tt.existingData, tt.errGet).
				Once()

			err := svc.Delete(tt.paramID)

			assert.Error(t, err)
			assert.EqualError(t, err, tt.messageError)

			repo.AssertExpectations(t)
		})
	}

	// result student after auto-increment ID
	expectedResult := model.Student{
		ID:   3,
		Name: "Rahman",
		Age:  20,
	}
	// service + mock repository
	svc, repo := newTestService()

	repo.On("GetAll").Return(service, nil).Once()

	// return nil setelah berhasil di tambah
	repo.On("SaveAll", mock.Anything).Return(nil).Once()

	result, err := svc.Create(createStudent)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	// semua mock terpenuhi
	repo.AssertExpectations(t)
}

// ======= UPDATED STUDENT
func TestStudentService_Updated(t *testing.T) {
	initial := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Siti", Age: 22},
	}
	svc, repo := newTestService()

	updateData := model.Student{
		Name: "Siti Updated",
		Age:  23,
	}

	repo.On("GetAll").Return(initial, nil).Once()
	repo.On("SaveAll", mock.Anything).Return(nil).Once()

	result, err := svc.Update(2, updateData)

	assert.NoError(t, err)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, "Siti Updated", result.Name)
	assert.Equal(t, 23, result.Age)

	repo.AssertExpectations(t)
}

// ======= GET ALL STUDENT
func TestStudentService_GetAll(t *testing.T) {
	initial := []model.Student{
		{ID: 1, Name: "Andi", Age: 21},
		{ID: 2, Name: "Siti", Age: 22},
	}
	svc, repo := newTestService()

	repo.On("GetAll").Return(initial, nil).Once()

	students, err := svc.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.Equal(t, 2, len(students))
	assert.Equal(t, "Andi", students[0].Name)
	assert.Equal(t, "Siti", students[1].Name)

	repo.AssertExpectations(t)
}

// ======= GET ALL STUDENT ERROR
func TestStudentService_GetAll_Error(t *testing.T) {
	svc, repo := newTestService() // instance service dan mock

	// ketika gagal return maka dia mengembalikan array dan error
	repo.On("GetAll").Return([]model.Student{}, utils.ErrFile).Once()

	students, err := svc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, students)
	assert.Equal(t, utils.ErrFile, err)

	repo.AssertExpectations(t)
}
