package main

import (
	"fmt"
	"sync"
)

func printOrderedWithChannels(m, n int) {
	channels := make([]chan int, m)
	for i := 0; i < m; i++ {
		channels[i] = make(chan int)
	}

	var wg sync.WaitGroup
	
	// 启动 m 个 goroutine
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 每个线程负责打印属于自己的那些数字
			// 线程 0 打印 1, 1+m, 1+2m...
			// 注意：这里为了简化逻辑，我们让每个线程循环尝试获取令牌
			// 更精确的做法是计算自己需要打印几次
			
			// 简单起见，我们让所有线程都跑起来，通过 channel 控制谁能动
			// 但实际上每个线程只需要处理 n/m 次左右
			
			current := id + 1
			for current < n {
				// 1. 等待自己的通道有信号 (令牌)
				<-channels[id]
				
				// 2. 打印
				fmt.Printf("%d ", current)
				
				// 3. 计算下一个数字
				current += m
				
				// 4. 将令牌传递给下一个线程 ( (id+1)%m )
				// 如果 current > n，其实不需要再传了，但为了代码简洁，
				// 我们可以传，只是接收方发现 current > n 后会退出
				nextID := (id + 1) % m
				channels[nextID] <- 1
			}
			if current == n {
				<-channels[id]
				fmt.Printf("%d ", current)
			}
		}(i)
	}

	// 启动第一个线程 (给 channel[0] 发送初始令牌)
	channels[0] <- 1

	wg.Wait()
	fmt.Println()
}

func main() {
	m := 3
	n := 10
	printOrderedWithChannels(m,n)
}