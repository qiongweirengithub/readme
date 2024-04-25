package gcore

import (
	"fmt"
	"time"

	gconnection "gserver-gateway.com/gconn"
)

func Helloworld() {
	fmt.Println("gateway center hello world")
}


func Start() {

	gcm := gconnection.GconManagers{}
	gcm.Init()

	gls := GListenerServer{}
	gls.Init(&gcm)
	gls.StartListener("tcp", "0.0.0.0:8888")

	// 控制main不退出
	holdmain()

}

func holdmain() {
	for range time.Tick(500 * time.Millisecond) {
		// fmt.Println("common" + TimeNow())
	}
}
