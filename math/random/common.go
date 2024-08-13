package random

import "time"

// ProbabilityTrue 在指定的百分比范围内返回true
// 0 <= percentage <= 100 超出范围将永远返回false
var seed = NewRandom(time.Now().Unix())

func ProbabilityTrue(percentage int) bool {
	if percentage < 0 || percentage > 100 {
		return false
	}
	randomNumber := seed.RandInt(100)
	return randomNumber < percentage
}
