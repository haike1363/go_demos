package main

import "testing"

func TestStudent_Nice(t *testing.T) {
	var stu Student
	stu.Nice()
}

func BenchmarkStudent_Nice(b *testing.B) {
	var stu Student
	for i := 0; i < 10; i++ {
		stu.Nice()
	}
}
