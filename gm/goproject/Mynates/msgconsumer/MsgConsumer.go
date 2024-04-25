package msgconsumer

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
	"gserver.com/mynatesclient"
)

func Helloworld() {
	fmt.Println("msg consumer hello world")
}

type Consumer struct {
	client *mynatesclient.MyClient
	topic  string
}

func (c *Consumer) Init(topic string, client *mynatesclient.MyClient) {
	c.client = client
	c.topic = topic
}

func (c *Consumer) Receive() {
	c.client.Receive(c.topic, func(msg *nats.Msg) {
		fmt.Println("consumer succ: " + string(msg.Data))
	})
}
