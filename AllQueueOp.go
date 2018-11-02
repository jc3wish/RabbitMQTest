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


	var qList *[]queueInfo
	if _,ok:=config["QueueList"];ok{
		qList = GetQueuesByConfig(config["QueueList"])
	}else{
		qList = GetQueuesByUrl(HttpUri,AmqpAdmin,AmqpPwd)
	}

	AllStartTime := time.Now().UnixNano() / 1e6

	OverCount := 0
	NeedWaitCount := 0
	ResultChan := make(chan int,NeedWaitCount)

	log.Println(key,"AllQueueOp start",AllStartTime)

	for keyI,qInfo := range *qList{
		m := make(map[string]string)
		if qInfo.Vhost == "/"{
			m["Uri"] = "amqp://"+AmqpAdmin+":"+AmqpPwd+"@"+AmqpUri+"/"
		}else{
			m["Uri"] = "amqp://"+AmqpAdmin+":"+AmqpPwd+"@"+AmqpUri+"/"+qInfo.Vhost
		}

		keyString := key+"-"+config["Method"]+strconv.Itoa(keyI)

		switch config["Method"] {
		case "all_write":
			m["ConnectCount"] = config["ConnectCount"]
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
			m["ConnectCount"] = config["ConnectCount"]
			m["QueueName"] = qInfo.Queue
			m["ConsumeTimeOut"] = config["ConsumeTimeOut"]
			if _,ok:=config["ConsumeCount"];ok{
				m["ConsumeCount"] = config["ConsumeCount"]
			}else{
				m["ConsumeCount"] = "0"
			}
			if _,ok:=config["AutoAck"];ok{
				m["AutoAck"] = config["AutoAck"]
			}else{
				m["AutoAck"] = "0"
			}
			if _,ok:=config["ExchangeName"];ok{
				m["ExchangeName"] = config["ExchangeName"]
			}else{
				m["ExchangeName"] = ""
			}
			NeedWaitCount++
			go SingleConsume(keyString,m,ResultChan)
			break
		default:
			m2 := make(map[string]string)
			NeedWaitCount += 2
			m2["Uri"] = m["Uri"]
			m2["QueueName"] = qInfo.Queue
			m2["ConnectCount"] = config["CosumeConnectCount"]
			m2["ConsumeTimeOut"] = config["ConsumeTimeOut"]
			if _,ok:=config["AutoAck"];ok{
				m2["AutoAck"] = config["AutoAck"]
			}else{
				m2["AutoAck"] = "0"
			}
			if _,ok:=config["ConsumeCount"];ok{
				m2["ConsumeCount"] = config["ConsumeCount"]
			}else{
				m2["ConsumeCount"] = "0"
			}
			go SingleConsume(keyString,m2,ResultChan)

			m["ConnectCount"] = config["WriteConnectCount"]
			m["DeliveryMode"] = config["DeliveryMode"]
			m["DateSize"] = config["DateSize"]
			m["ChannelCount"] = config["ChannelCount"]
			m["ChanneWriteCount"] = config["ChanneWriteCount"]
			m["WaitConfirm"] = config["WaitConfirm"]
			m["WriteTimeOut"] = config["WriteTimeOut"]
			if _,ok:=config["ExchangeName"];ok{
				m["ExchangeName"] = config["ExchangeName"]
			}else{
				m["ExchangeName"] = ""
			}
			m["RoutingKey"] = qInfo.Queue
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
	AllEndTime := time.Now().UnixNano() / 1e6
	log.Println(key,"AllQueueOp end",AllEndTime," time(ms):",AllEndTime-AllStartTime)
}
