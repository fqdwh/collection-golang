# 基础语法

## Task1

```go
// 2.go
fmt.Fscan(os.Stdin, &h) // 读取位置，地址
fmt.Scan(&h)            // 只需要地址，默认从 os.Stdin 读
```
* 两种方式都将空格/换行当做分隔符
---
```go
// 6.go
file,err:= os.Create("filename")
if err!=nil{

}
defer file.Close()
fmt.Fprintf(file,"str")
```
* defer 在错误检查后关闭文件，否则 Close nil 指针会 panic
* os.Create() 文件存在覆盖，不存在创建
---
切片和数组的区别

1. 切片可以自动扩容，数组长度在创建时就已经固定
2. 不同大小数组不兼容，属于不同类型
3. 切片可以通过数组初始化
4. 切片可以为空
5. 切片可以使用 len cap append copy 等函数
```go
// 切片创建方法
var s []int
s := []int
s := []int{...}
s := make([]int,len,cap)
s := arr[...]
s := s1[...]

// map创建方法
m := make(map[keytype]valuetype,cap)
m := map[string]int{
    "apple":1,
}
m := map[string]int{}   // 空 map 而不是 nil， 可以安全写入
var m map[string]int    // m==nil
```
---
Bouns3
```go
package main

import (
	"fmt"
)

// 生成器
// 不断产生整数列
func generate(ch chan int) {
	for i := 2; ; i++ {     // 死循环
		ch <- i             // 生成数字发送到通道
	}
}

// 过滤器
// 过滤能被质数整除的数
func filter(in chan int, out chan int, prime int) {
	for {
		num := <-in         // 输入从通道接收数字
		if num%prime != 0 {
			out <- num      // 不能被当前质数整除就发送到输出通道
		}
	}
}

func main() {
	ch := make(chan int)
	go generate(ch)
	for i := 0; i < 6; i++ {
		prime := <-ch 
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime)
		ch = out
	}
	// 第一轮 prime = 2 
	// ch 里筛出 2 的倍数，其余进入 out 
	// ch 指向 out
	// 在筛出 2 的倍数的基础下再筛出 3 的倍数
	// 循环
}
```
1. 实现了一个基于并发流水线的质数筛选器
2. 采取 Goroutines 并发运行，通过 channels 串联进程
3. 没有性能上的提升：每次传递的通信成本远大于计算成本；随着质数越来越多，会创建越来越多的 goroutine 和 channel，造成资源浪费

---



# OSPP

* 爬虫 AI 部分：
1. 响应 JSON 解析到结构体
2. mysql 初始化以及插入新列


***爬取速度***


```
// 无 PDF
// 完全单程  1m40.1237254s
// 列表单程 详细信息最多 50 并行 50.3590556s
// 列表 12 并行 详细信息 50 并行 5.7801475s
// 加速比 ≈ 1:17.3

// 有 PDF
// 列表 12 并行 详细信息 50 并行 PDF不并行 1m5.1900843s
```
