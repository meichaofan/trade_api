package common_test

import (
	"fmt"
	"testing"
	"trade_api/src/main/web/cli/common"
)

func TestCalRate(t *testing.T) {
	rate := common.CalRate("cny")
	fmt.Println(rate)
}
