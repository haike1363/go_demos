package main

import (
	"bytes"
	"fmt"
)

func parsePlaceHolder(s string) ([]string, error) {
	var result []string
	match := bytes.Buffer{}
	state := 0
	index := 0
	for i, c := range s {
		if state == 0 {
			if c == '$' {
				state = 1
				match.WriteString("$")
				index = i
			}
		} else if state == 1 {
			if c == '{' {
				state = 2
				match.WriteString("{")
			} else {
				state = 0
				match = bytes.Buffer{}
			}
		} else if state == 2 {
			match.WriteRune(c)
			if c == '}' {
				state = 0
				result = append(result, match.String())
				match = bytes.Buffer{}
			}
		}
	}
	if state != 0 {
		return nil, fmt.Errorf("bad parsePlaceHolder failed %v index %v miss match %v", s, index, match.String())
	}
	return result, nil
}

func main() {

	match, err := parsePlaceHolder("$abcbcd")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range match {
		fmt.Println(v)
	}
}
