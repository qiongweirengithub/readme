package common

import (
	"fmt"
	"strconv"
	"time"
)

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Tick() {
	for range time.Tick(500 * time.Millisecond) {
		fmt.Println("common" + TimeNow())
	}
}
func Holdmain() {
	for range time.Tick(500 * time.Millisecond) {
		// fmt.Println("common" + TimeNow())
	}
}

func Tick2(timec int) {
	timer := time.NewTimer(time.Duration(timec) * time.Second)
	str1 := strconv.Itoa(timec)
	for {
		timer.Reset(time.Duration(timec) * time.Second) // 这样来复用 timer 和修改执行时间
		select {
		case <-timer.C:
			fmt.Println("每隔" + str1 + "秒执行任务")
		}
	}
}
