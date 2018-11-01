package RabbitMQTest

import (
	"github.com/streadway/amqp"
	"log"
)

type Channel struct {
	ch *amqp.Channel
	confirmWait chan amqp.Confirmation
	WriteTimeOut int
	ConsumeTimeOut int
}

func(Channel *Channel) SetWriteTimeOut(WriteTimeOut int){
	Channel.WriteTimeOut = WriteTimeOut
}

func(Channel *Channel) SetConsumeTimeOut(ConsumeTimeOut int){
	Channel.ConsumeTimeOut = ConsumeTimeOut
}

type Conn struct {
	Uri    string
	status string
	conn   *amqp.Connection
	err    error
	expir  string
	chList []*Channel
}

func NewConn(Uri string) *Conn{
	c := &Conn{Uri:Uri,chList:make([]*Channel,0)}
	c.conn, c.err = amqp.Dial(Uri)
	if c.err != nil{
		log.Println(Uri,"connect err:",c.err)
	}
	return c
}

func(This *Conn)NewChannel(wait bool) (*Channel,error) {
	ch, err := This.conn.Channel()
	if err != nil {
		return nil, err
	}
	if wait == true {
		ch.Confirm(false)
		confirmWait := make(chan amqp.Confirmation, 1)
		ch.NotifyPublish(confirmWait)
		return &Channel{ch: ch, confirmWait: confirmWait}, nil
	}else{
		return &Channel{ch: ch, confirmWait: nil}, nil
	}
}