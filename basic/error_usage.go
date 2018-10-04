package main

type PathError struct {
	Op string
}

func (self *PathError) Error() string {
	return self.Op
}

func main() {

}
