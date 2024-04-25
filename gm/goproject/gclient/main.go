package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client starting")

	con, err := net.Dial("tcp", "192.168.0.100:8888")
	if err != nil {
		fmt.Println("connect to server fail", err)
		panic(err)
	}
	establish(con)

	go tickWrite(con)
	go read(con)

	holdmain()

}

func holdmain() {
	for range time.Tick(500 * time.Millisecond) {
		// fmt.Println("common" + TimeNow())
	}
}

func establish(con net.Conn) {
	var data = make(map[string]interface{})
	data["type"] = "establish"
	data["userid"] = "test_userid001"
	a := mapToJson(data)
	_, err := con.Write([]byte(a))
	if err != nil {
		fmt.Println("data send error", err)
		panic(err)
	}
	time.Sleep(2000)
}

func read(con net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := con.Read(buffer)
		if err != nil {
			fmt.Println("conn read error", err)
			panic(err)
		}
		msg := string(buffer[0:n])
		fmt.Println("client msg-:" + msg)
	}
}

func tickWrite(con net.Conn) {
	var data = make(map[string]interface{})
	data["type"] = "hello"
	data["userid"] = "test_userid001"
	a := mapToJson(data)
	_, err := con.Write([]byte(a))
	if err != nil {
		fmt.Println("data send error", err)
		panic(err)
	}
}

func timeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func mapToJson(param map[string]interface{}) string {
	//json转map
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func jsonToMap(str string) map[string]interface{} {
	//map 转json

	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}
