package main

import (
	"encoding/base64"
	"io/ioutil"
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"time"
	"github.com/streadway/amqp"

)



func main(){
	doMQTest()
	//httpTest()
}

func httpTest(){
	type queueInfo struct {
		Vhost string `json:"vhost"`
		Queue string `json:"name"`
	}
	var qList []queueInfo

	qList = make([]queueInfo,0)

	content := GetQueues("http://10.4.4.199:15672/api/queues","admin","admin")
	json.Unmarshal([]byte(content),&qList)
	for key,v := range qList{
		log.Println(key,v.Vhost,v.Queue)
	}
}

func GetQueues(url string,user string,pwd string) string{
	req, _ := http.NewRequest("GET", url, nil)

	//req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Authorization","Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pwd)))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
	return string(body)
}

func doMQTest(){
	uri := "amqp://admin:admin@10.4.4.199:5672/testvhost"
	exchangeName:="amq.direct"
	exchangeName="amq.default"
	exchangeName=""
	queueName:="testQueue"
	conn := NewConn(uri)
	ch,err := conn.NewChannel(true)
	if err!=nil{
		log.Println(err)
		return
	}
	ch.SetWriteTimeOut(5)
	var Mode uint8 = 1
	var Body []byte = []byte("ccc")
	_,err2 := SendMQ(ch,&exchangeName,&queueName,&Mode,&Body)
	if err2 != nil{
		log.Println(err)
	}
}

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

func SendMQ(Channel *Channel,exchange *string,routingkey *string,DeliveryMode *uint8,c *[]byte)(bool,error){
	err := Channel.ch.Publish(
		*exchange,     // exchange
		*routingkey, // routing key
		true,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:   *c,
			DeliveryMode:	*DeliveryMode,
			//Expiration:	This.expir,
		})
	if err != nil{
		return false,err
	}
	if Channel.confirmWait != nil{
		select {
		case d := <-Channel.confirmWait:
			if d.DeliveryTag >= 0{
				return true,nil
			}
			return false,fmt.Errorf("unkonw err")
			break
		case <-time.After(time.Duration(Channel.WriteTimeOut) * time.Second):
			return false,fmt.Errorf("server no response")
			break
		}
	}else{
		return true,nil
	}
	return false,fmt.Errorf("unkonw err")
}
