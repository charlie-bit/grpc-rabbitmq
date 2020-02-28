package utils

import "log"

/**
	Author:charlie
	Description: 输出错误信息
	Time:2020-2-24
*/
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
