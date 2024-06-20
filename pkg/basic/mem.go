package main

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

var s = ""

func init() {
	var buf bytes.Buffer
	buf.Grow(1024 * 1024 * 6)
	for i := 0; i < 64*1024*6; i++ {
		buf.WriteString("1234567890123456")
	}
	s = buf.String()
}

func getBufferString() string {
	var buf bytes.Buffer
	for i := 0; i < 2; i++ {
		buf.WriteString(s)
	}
	return buf.String()
}

type Node struct {
	index int
}
type Selection struct {
	nodes []*Node
}

func newSelection(node *Node) *Selection {
	s := &Selection{}
	s.nodes = append(s.nodes, node)
	return s
}

func (this *Selection) Each(f func(int, *Selection)) {
	for i, n := range this.nodes {
		f(i, newSelection(n))
	}
}

func main() {
	s := &Selection{}
	for i := 0; i < 10000; i++ {
		s.nodes = append(s.nodes, &Node{index: i})
	}
	s.Each(func(i int, s *Selection) {
		ss := getBufferString()
		println(len(ss))
		if i%20 == 0 {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			fmt.Printf("已使用内存：%v %v\n", memStats.Alloc, memStats.TotalAlloc)
		}
		time.Sleep(1 * time.Second)
	})
}
