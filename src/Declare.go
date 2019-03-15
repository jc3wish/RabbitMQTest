package RabbitMQTest

import (
	"time"
	"strings"
	"strconv"
	"log"
	"github.com/streadway/amqp"
)

type declareQueueInfo struct {
	Name string
	Durable bool
	AutoDelete bool
}


func GetDeclareQueuesByConfig(QueueList string) []declareQueueInfo{
	s := strings.Split(QueueList,",")
	qList := make([]declareQueueInfo,0)
	for _,v := range s{
		q := strings.Split(v,":")
		var Name string
		var Durable bool = true
		var AutoDelete bool = false
		n := len(q)
		if n > 0{
			Name = q[0]
		}else{
			continue
		}
		if n > 1 && strings.Trim(q[1],"") != "true"{
			Durable = false
		}

		if n > 2 && strings.Trim(q[2],"") == "true"{
			AutoDelete = true
		}

		qList = append(qList, declareQueueInfo{Name,Durable,AutoDelete})
	}
	return qList
}

func GetDeclareQueuesByConfig2(prefix string,DurableString string,AutoDeleteString string,Count int) []declareQueueInfo{
	var Durable bool = true
	var AutoDelete bool = false
	if DurableString != "true"{
		Durable = false
	}
	if AutoDeleteString == "true"{
		AutoDelete = true
	}

	qList := make([]declareQueueInfo,0)

	for i:=0;i<Count;i++{
		qList = append(qList, declareQueueInfo{prefix+strconv.Itoa(i),Durable,AutoDelete})
	}

	return qList

}

func Only_Declare(key string,config map[string]string,resultDataChan chan *Result){
	ResultData := NewResult()
	AmqpUri := config["Uri"]

	var qList []declareQueueInfo
	if _,ok:=config["QueueList"];ok{
		qList = GetDeclareQueuesByConfig(config["QueueList"])
	}else{
		qList = GetDeclareQueuesByConfig2(config["QueuePrefix"],config["QueueDurable"],config["QueueAutoDelete"],GetIntDefault(config["QueueCount"],0))
	}

	var (
		ExchangeName string = "amq.direct"
		ExchangeType string = "direct"
		ExchangeDurable bool = true
		ExchangeAutoDelete bool = false
		RoutingKey string = ""
	)

	if config["ExchangeName"] != ""{
		ExchangeName = config["ExchangeName"]
	}
	if config["ExchangeType"] != ""{
		ExchangeType = config["ExchangeType"]
	}

	if config["ExchangeDurable"] != "" && config["ExchangeDurable"] != "true"{
		ExchangeDurable = false
	}

	if config["ExchangeAutoDelete"] == "true"{
		ExchangeAutoDelete = true
	}

	RoutingKey = config["RoutingKey"]

	AllStartTime := time.Now().UnixNano() / 1e6
	Conn := NewConn(AmqpUri)
	Channel,err := Conn.NewChannel(false)
	if err != nil{
		log.Println(key,err)
		return
	}
	T := make(amqp.Table,0)
	Channel.ch.ExchangeDeclare(ExchangeName,ExchangeType,ExchangeDurable,ExchangeAutoDelete,false,false,T)
	for _,Q := range qList{
		Channel.ch.QueueDeclare(Q.Name,Q.Durable,Q.AutoDelete,false,false,T)
		if RoutingKey == ""{
			Channel.ch.QueueBind(Q.Name,Q.Name,ExchangeName,false,T)
		}else{
			Channel.ch.QueueBind(Q.Name,RoutingKey,ExchangeName,false,T)
		}
	}
	AllEndTime := time.Now().UnixNano() / 1e6
	UseTime := float64(AllEndTime-AllStartTime)
	log.Println(key,"Declare Count:",len(qList),"endTime:",AllEndTime," had use time(ms):",UseTime)
	resultDataChan <- ResultData
}
