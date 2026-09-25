package main

import( 
	"fmt"
)
func judge(n int) bool {
	if n%100 == 0 {
		if n%400 == 0 {
			return true
		} else {
			return false
		}
	} else {
		if n%4 == 0 {
			return true
		} else {
			return false
		}
	}
}

func main() {
	var a,b int
	num := 0
	fmt.Scanf("%d%d",&a,&b)
	var l [400]int
	for i := a;i<=b;i++ {
		if judge(i) {
			l[num] = i
			num++
	    }
    }
	fmt.Println(num)
	for i := 0;i<num;i++ {
		fmt.Print(l[i]," ")
	}
}