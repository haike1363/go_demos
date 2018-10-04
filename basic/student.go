package basic

type People interface {
}

type Student struct {
    Name string
    Year int
}

//func (*Student)Set(self (*Student), Name string, Year int) string {
//    return Name
//}