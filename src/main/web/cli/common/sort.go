package common

import "trade_api/src/main/web/cli/data"

//这里数据量不大，就写个冒泡吧
func Sort(amountList []*data.ExchangeAmount, order string) {
	length := len(amountList)
	for i := 0; i < length; i++ {
		for j := 0; j+1 < length-i; j++ {
			if order == "asc" && amountList[j].TotalUsd > amountList[j+1].TotalUsd {
				amountList[j], amountList[j+1] = amountList[j+1], amountList[j]
			} else if order == "desc" && amountList[j].TotalUsd < amountList[j+1].TotalUsd {
				amountList[j], amountList[j+1] = amountList[j+1], amountList[j]
			}
		}
	}
}

func MySort(nums []int, order string) {
	length := len(nums)
	for i := 0; i < length; i++ {
		for j := 0; j+1 < length-i; j++ {
			if order == "asc" && nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			} else if order == "desc" && nums[j] < nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}
}
