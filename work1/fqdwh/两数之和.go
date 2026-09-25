package main

import(
	"fmt"
)

func twoSum(nums []int, target int) []int {
	m := map[int]int{}
	for index,val := range(nums) {
		_,ok := m[target-nums[index]]
		if ok {
			return []int{m[target-nums[index]],index}
		}
		m[val] = index
	}
	return []int{}
}
func main() {
	nums := []int{2,7,1,15}
	var target int = 9
	s := twoSum(nums,target)
	fmt.Println(s)
}