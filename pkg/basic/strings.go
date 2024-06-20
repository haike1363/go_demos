package main

import (
	"bytes"
	"fmt"
	"strings"
)

func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func OptimizeK(v int64) string {
	tmp := ReverseString(fmt.Sprint(v))
	var buff = bytes.Buffer{}
	for i, c := range tmp {
		if i != 0 && i%3 == 0 {
			buff.WriteRune(',')
		}
		buff.WriteRune(c)
	}
	return ReverseString(buff.String())
}
func main() {
	tmp := strings.Split("createTime", ":")
	fmt.Println(len(tmp))
	fmt.Println(tmp)
	fmt.Println(tmp[0])
	fmt.Println(strings.Join(tmp[1:], ":"))

	fmt.Println(OptimizeK(0))
	fmt.Println(OptimizeK(1))
	fmt.Println(OptimizeK(10))
	fmt.Println(OptimizeK(100))
	fmt.Println(OptimizeK(1000))
	fmt.Println(OptimizeK(100000000))
	{
		values := strings.Split("instanceid,id", ",")
		fmt.Println(len(values))
		fmt.Println(values)
	}
	values := strings.Split("~456", "~")
	fmt.Println(len(values))
	fmt.Println(values)

	values = strings.Split("456~", "~")
	fmt.Println(len(values))
	fmt.Println(values)

	values = strings.Split("123~456", "~")
	fmt.Println(len(values))
	fmt.Println(values)

	values = strings.Split("~", "~")
	fmt.Println(len(values))
	fmt.Println(values)

	s := "jdbc:mysql://172.17.80.46:3306/hivemetastore?useSSL=false&amp;createDatabaseIfNotExist=true&amp;characterEncoding=UTF-8"
	parts := strings.Split(
		strings.Replace(
			strings.Replace(
				strings.Replace(s, "//", "", -1), "/", ":", -1), "?", ":", -1), ":")
	for i, p := range parts {
		fmt.Println(i, p)
	}
}
