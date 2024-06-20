package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"syscall"
)

func main() {
	var stat syscall.Statfs_t
	syscall.Statfs("/Users/easechen/tmp/Hive.tar.gz", &stat)
	fmt.Println(stat.Bsize)

	// 打开 tar.gz 文件
	file, err := os.Open("/Users/easechen/tmp/Hive.tar.gz")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建 gzip.Reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	// 创建 tar.Reader
	tarReader := tar.NewReader(gzipReader)

	// 遍历 tar 文件并计算文件大小
	var totalSize int64
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if header.Typeflag == tar.TypeReg {
			// 物理磁盘需要 4KB对齐
			totalSize += int64(math.Ceil(float64(header.Size)/4096) * 4096)
		}
	}

	// 输出解压缩后的文件大小
	fmt.Printf("解压缩后的文件大小为 %d KB\n", totalSize/1024)
}
