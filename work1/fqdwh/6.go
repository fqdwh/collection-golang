package main

import(
	"fmt"
	"os"
)

func main() {
	file,err := os.Create("ninenine.txt") // 不存在创建，存在清空
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer file.Close()

	for i:=1;i<=9;i++ {
		for j:=1;j<=i;j++ {
			fmt.Fprintf(file," %d * %d = %2d ",j,i,i*j)
		}
		fmt.Fprintf(file,"\n")
	}
}