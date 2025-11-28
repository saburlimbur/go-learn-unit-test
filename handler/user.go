package handler

import (
	"fmt"
	"session-9/model"
	"session-9/service"
	"strings"
)

type StudentHandler struct {
	Svc service.StudentServiceInterface
}

func NewStudentHandler(svc service.StudentServiceInterface) *StudentHandler {
	return &StudentHandler{Svc: svc}
}

// ListStudents returns formatted list of students
func (h *StudentHandler) ListStudents() (string, error) {
	students, err := h.Svc.GetAll()
	if err != nil {
		return "", err
	}

	if len(students) == 0 {
		return "No students found.\n", nil
	}

	var b strings.Builder
	b.WriteString("Students:\n")
	for _, st := range students {
		// format: ID - Name (Age)
		fmt.Fprintf(&b, "- %d: %s (%d)\n", st.ID, st.Name, st.Age)
	}
	return b.String(), nil
}

// CreateStudent handles creating a new student and returns message
func (h *StudentHandler) CreateStudent(name string, age int) (string, error) {
	input := model.Student{
		Name: name,
		Age:  age,
	}

	created, err := h.Svc.Create(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Student created: ID=%d, Name=%s, Age=%d\n",
		created.ID, created.Name, created.Age), nil
}
