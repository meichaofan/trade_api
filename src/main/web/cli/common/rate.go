package common

import (
	"github.com/tidwall/gjson"
	"strings"
)

/**
各国货币和美元汇率
*/

func init() {
	CnyUsdRate = CalRate("cny")
}

var CnyUsdRate float64

const (
	//各国货币和美元之间汇率
	RateApi = "http://49.51.46.37:14000/v1/getexchange"
)

// 1美元 == xxx XXX
func CalRate(base string) float64 {
	symbol := strings.ToUpper("USD" + base)
	url := RateApi
	content := HttpGet(url)
	ret := gjson.ParseBytes(content)
	status := ret.Get("status.status_code").Int()
	if status == 200 {
		rates := ret.Get("data.quotes").Map()
		if rate, exist := rates[symbol]; exist {
			return rate.Float()
		} else {
			panic(base + "is not exist")
		}
	}
	return 0
}
