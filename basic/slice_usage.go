package main

import "fmt"

// 此处不能传数组，数组是值传递
func modifyArray(slice []int) {
    slice[0] = 2001
}

func main() {
    var myArray [10]int = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    var mySlice []int = myArray[:]
    // 传递切片修改数组
    modifyArray(mySlice)
    fmt.Println(myArray)



    // 定义Slice设置初始值以及最大容量
    mySlice5_max10 := make([]int, 5, 10)

    fmt.Println(cap(mySlice5_max10), len(mySlice5_max10))
    mySlice5_max10 = append(mySlice5_max10, 1, 2, 3)

    // 切片就是引用效果
    newSlice := mySlice5_max10[0:3]
    newSlice[0] = 1998
    fmt.Println(newSlice)
    fmt.Println(mySlice5_max10)

    mySliceValus := []int{1, 2, 3}
    // 自动扩展容量，如果超过最大容量，会生成一个新切片，老切片内容不变，而之前有N切片从老切片生成，那么N切片的修改不会作用于新生成的切片
    // 如果合适，则直接修改原来的切片，如果之前有N切片从老切片生成，那么N切片的修改继续作用
    mySlice5_max15 := append(mySlice5_max10, mySliceValus...)
    newSlice[1] = 1999
    fmt.Println(mySlice5_max10)
    fmt.Println(mySlice5_max15)
    // 拷贝
    copy(mySlice5_max15, mySlice5_max10)
}
