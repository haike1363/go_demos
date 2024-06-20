package main

import (
    "encoding/json"
    "fmt"
)

func toJson(v interface{}) string {
    b, err := json.MarshalIndent(v, "", "\t")
    if err != nil {
        return fmt.Sprintf("%v toJsonStringErr %v", v, err)
    }
    return string(b)
}

func main() {
    tsMap := map[string]interface{}{}
    tsMap["a"] = nil
    tsMap["b"] = 1
    fmt.Println(toJson(tsMap))
    var resultList = []string{}
    fmt.Print(toJson(resultList), "\n")

    var resultMap = map[string]interface{}{}
    fmt.Print(toJson(resultMap),"\n")
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
