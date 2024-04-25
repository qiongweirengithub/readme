package msgproducer

import (
	"fmt"
	"time"

	"gserver.com/mynatesclient"
	"juan.com/common"
)

func Helloworld() {
	fmt.Println("msg producer hello world")
}

type Producer struct {
	client *mynatesclient.MyClient
	topic  string
}

func (p *Producer) Init(topic string, client *mynatesclient.MyClient) {
	p.client = client
	p.topic = topic
}

func (p *Producer) Send(data string) {
	p.client.Send(p.topic, data)
}

func (p *Producer) TickSend(data string) {

	for range time.Tick(2000 * time.Millisecond) {
		fmt.Println(p.topic + ",sending: " + data + common.TimeNow())
		p.client.Send(p.topic, data+common.TimeNow())
	}
}
