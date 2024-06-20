package main

import (
	"testing"
)

func TestStudent_Nice(t *testing.T) {
	type fields struct {
		Name string
		Year int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &Student{
				Name: tt.fields.Name,
				Year: tt.fields.Year,
			}
			self.Nice()
		})
	}
}

func TestStudent_Say(t *testing.T) {
	type fields struct {
		Name string
		Year int
	}
	type args struct {
		a int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &Student{
				Name: tt.fields.Name,
				Year: tt.fields.Year,
			}
			self.Say(tt.args.a)
		})
	}
}
