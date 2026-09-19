package main

import (
	"fmt"
	"slices"
)

type Student struct {
	Name  string
	Score float64
}

func printStudents(students []Student) {
	for _, student := range students {
		fmt.Printf("Name: %s; Score: %.1f\n", student.Name, student.Score)
	}
}

func addStudent(students []Student, student Student) []Student {
	return append(students, student)
}

func removeStudentByName(students []Student, name string) []Student {
	for i, student := range students {
		if student.Name == name {
			return append(students[:i], students[i+1:]...)
		}
	}
	return students
}

func sortByScore(students []Student) {
	slices.SortFunc(students, func(a, b Student) int {
		if a.Score > b.Score {
			return -1
		}
		if a.Score < b.Score {
			return 1
		}
		return 0
	})
}
func main() {
	students := []Student{
		{Name: "An", Score: 9.4},
		{Name: "Tuan", Score: 7.4},
		{Name: "Long", Score: 4.2},
	}
	fmt.Println("Danh sach ban dau: ")
	printStudents(students)
	students = addStudent(students, Student{
		Name:  "Khang",
		Score: 10,
	})
	fmt.Println("Sau khi them Khang: ")
	printStudents(students)
	sortByScore(students)
	fmt.Println("Sau khi sap xep:")
	printStudents(students)
	copyStudents := make([]Student, len(students))
	copy(copyStudents, students)
	copyStudents = removeStudentByName(copyStudents, "Tuan")
	fmt.Println("Ban sao sau khi xoa Tuan:")
	printStudents(copyStudents)
	fmt.Println("Danh sach goc:")
	printStudents(students)

}
