package RabbitMQTest

import (
	"time"
	"log"
	"strconv"
)

func AllQueueOp(key string,config map[string]string,resultBackChan chan int){
	defer func() {
		resultBackChan <- 1
	}()

	AmqpUri := config["AmqpUri"]
	HttpUri := config["HttpUri"]
	AmqpAdmin := config["AmqpAdmin"]
	AmqpPwd := config["AmqpPwd"]

	qList := GetQueuesByUrl(HttpUri,AmqpAdmin,AmqpPwd)

	AllStartTime := time.Now().Unix()

	OverCount := 0
	NeedWaitCount := 0
	ResultChan := make(chan int,NeedWaitCount)

	log.Println(key,"AllQueueOp start",AllStartTime)

	for keyI,qInfo := range *qList{
		m := make(map[string]string)
		m["ConnectCount"] = config["ConnectCount"]
		if qInfo.Vhost == "/"{
			m["Uri"] = "amqp://"+AmqpAdmin+":"+AmqpPwd+"@"+AmqpUri+"/"
		}else{
			m["Uri"] = "amqp://"+AmqpAdmin+":"+AmqpPwd+"@"+AmqpUri+"/"+qInfo.Vhost
		}

		keyString := key+"-"+config["Method"]+strconv.Itoa(keyI)

		switch config["Method"] {
		case "all_write":
			m["DeliveryMode"] = config["DeliveryMode"]
			m["DateSize"] = config["DateSize"]
			m["ChannelCount"] = config["ChannelCount"]
			m["ChanneWriteCount"] = config["ChanneWriteCount"]
			m["WaitConfirm"] = config["WaitConfirm"]
			m["WriteTimeOut"] = config["WriteTimeOut"]
			m["ExchangeName"] = ""
			m["RoutingKey"] = qInfo.Queue
			go SingleSend(keyString,m,ResultChan)
			NeedWaitCount++
			break
		case "all_consume":
			m["QueueName"] = qInfo.Queue
			m["ConsumeTimeOut"] = config["ConsumeTimeOut"]
			NeedWaitCount++
			go SingleConsume(keyString,m,ResultChan)
			break
		default:
			NeedWaitCount += 2
			m["QueueName"] = qInfo.Queue
			m["ConnectCount"] = config["CosumeConnectCount"]
			m["ConsumeTimeOut"] = config["ConsumeTimeOut"]
			go SingleConsume(keyString,m,ResultChan)

			m["DeliveryMode"] = config["DeliveryMode"]
			m["DateSize"] = config["DateSize"]
			m["ChannelCount"] = config["ChannelCount"]
			m["ChanneWriteCount"] = config["ChanneWriteCount"]
			m["WaitConfirm"] = config["WaitConfirm"]
			m["WriteTimeOut"] = config["WriteTimeOut"]
			m["ExchangeName"] = ""
			m["RoutingKey"] = qInfo.Queue
			m["ConnectCount"] = config["WriteConnectCount"]
			go SingleSend(keyString,m,ResultChan)
			break
		}
	}

	if NeedWaitCount == 0{
		return
	}
	for{
		<- ResultChan
		OverCount++
		if OverCount >= NeedWaitCount{
			break
		}
	}
	AllEndTime := time.Now().Unix()
	log.Println(key,"AllQueueOp end",AllEndTime," time(ms):",AllEndTime-AllStartTime)
}
