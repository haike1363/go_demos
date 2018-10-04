package main

import "fmt"

func main() {
    // 指定容量
    myMap := make(map[string]int, 3)
    myMap["A"] = 1
    myMap["B"] = 2

    // 查找
    it, ok := myMap["A"]
    if ok {
        fmt.Println("value ", it)
    } else {
        fmt.Println("A not exist")
    }

    // 定义并初始化
    myMap = map[string]int{
        "A": 1,
        "B": 2,
        "C": 3,
    }
    // 删除元素
    delete(myMap, "B")
    fmt.Println(myMap)
}
