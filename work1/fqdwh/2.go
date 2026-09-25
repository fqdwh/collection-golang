package main

import( 
	"fmt"
	"os"
)

func main() {
	l := make([]int,10)
	for i:=0;i<10;i++ {
		fmt.Fscan(os.Stdin,&l[i])
	}
	var h int
	fmt.Fscan(os.Stdin,&h)
	h+=30
	num := 0;
	for i:=0;i<10;i++ {
		if l[i] <= h{
			num++
		}
	}
	fmt.Println(num)
}