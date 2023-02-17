package utils_test

import (
	"testing"
	"time"
	"mufe_service/camp/utils"
)

func TestUseTime(t *testing.T) {
	var u = utils.UseTime("get")
	time.Sleep(1 * time.Second)
	u.Mark("111111111")
	time.Sleep(2 * time.Second)
	u.Mark("2222")
	time.Sleep(3 * time.Second)
	u.Mark("3333")
}
