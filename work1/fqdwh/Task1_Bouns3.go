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