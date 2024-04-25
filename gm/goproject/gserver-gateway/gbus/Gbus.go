package gbus

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
)

func Helloworld() {
	fmt.Println("gbus hello world")
}

type GBus struct {
	client *MyClient
}

func Init(nates_server string) GBus {
	client := &MyClient{}
	client.Init(nates_server)

	msgBus := GBus{}
	msgBus.client = client
	return msgBus
}

func (c *GBus) Receive(topic string, grp string, cb nats.MsgHandler) (* nats.Subscription, error)  {
	nc, err := c.client.cn.Subscribe(topic, cb)
	if err != nil {
		fmt.Println("err", err)
		return nil,err
	}
	return nc, nil
	// c.client.cn.QueueSubscribe(topic, grp, cb)
}

func (c *GBus) UnReceive(topic string, grp string, cb nats.MsgHandler) {
	// c.client.cn.Subscribe(topic, cb)

	chin, err := c.client.cn.Subscribe(topic, cb)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	chin.Unsubscribe()
}

func (c *GBus) Send(topic string, data string) {
	err := c.client.cn.Publish(topic, []byte(data))
	if err != nil {
		fmt.Println("publish err:"+data, err)
	}
}
