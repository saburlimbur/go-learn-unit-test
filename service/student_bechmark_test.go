package service

import (
	"session-9/model"
	"session-9/repository"
	"testing"
)

func newBenchmarkService() (*StudentService, *repository.MockStudentRepository) {
	repo := &repository.MockStudentRepository{}
	svc := NewStudentService(repo)
	return svc, repo
}

// func BenchmarkStudentService_Create(b *testing.B) {
// 	svc, _ := newBenchmarkService()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		_, err := svc.Create(model.Student{
// 			Name: "User",
// 			Age:  20,
// 		})
// 		if err != nil {
// 			b.Fatalf("unexpected error: %v", err)
// 		}
// 	}
// }

func BenchmarkStudentService_GetByID(b *testing.B) {
	initial := []model.Student{
		{ID: 1, Name: "A", Age: 21},
		{ID: 2, Name: "B", Age: 22},
		{ID: 3, Name: "C", Age: 23},
	}

	svc, repo := newBenchmarkService()
	repo.
		On("GetAll").
		Return(initial, nil).
		Times(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.GetByID(3)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkStudentService_GetAll(b *testing.B) {
	var students []model.Student
	for i := 1; i <= 100; i++ {
		students = append(students, model.Student{
			ID:   i,
			Name: "Student",
			Age:  20,
		})
	}

	svc, repo := newBenchmarkService()
	repo.
		On("GetAll").
		Return(students, nil).
		Times(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := svc.GetAll()
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
