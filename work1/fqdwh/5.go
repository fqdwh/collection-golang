package main

import(
	"fmt"
)

func main() {
	l := make([]int,0,50)
	for i:=1;i<=50;i++ {
		l = append(l, i)
	}
	num := 0
	for i:= 1;i<=50;i++ {
		if i%3 != 0 {
			l[num] = i
			num++
		}
	}
	l = l[:num]
	l = append(l,114514)
	fmt.Println(l)
}