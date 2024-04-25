package main

import (
	"fmt"

	"gserver.com/msgconsumer"
	"gserver.com/msgproducer"
	"gserver.com/mynatesclient"
	"juan.com/common"
)

func main() {
	fmt.Println(common.TimeNow())
	// common.Tick()
	// common.Tick2(1)
	mynatesclient.Helloworld()
	msgconsumer.Helloworld()
	msgproducer.Helloworld()

	client := &mynatesclient.MyClient{}
	client.Helloworld()
	client.Init("nats://db_nats:4333,nats://qw_nats:4333")

	topic := "hello001"
	p := &msgproducer.Producer{}
	p.Init(topic, client)
	c := &msgconsumer.Consumer{}
	c.Init(topic, client)

	go c.Receive()
	go p.TickSend("hello client")

	common.Holdmain()
}
