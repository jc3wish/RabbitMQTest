package RabbitMQTest

import (
	"log"
	"time"
)

func OnlyConnect(key string,config map[string]string,resultBackChan chan int){
	AmqpUri := config["Uri"]
	ConnectCount := GetIntDefault(config["ConnectCount"],1)
	ChannelCount := GetIntDefault(config["ChannelCount"],0)
	ConnectTimeOut := GetIntDefault(config["ConnectTimeOut"],100)
	connList := make([]*Conn,ConnectCount)
	for i:=0;i<ConnectCount;i++{
		connList[i] = NewConn(AmqpUri)
		if connList[i].err != nil{
			log.Println(AmqpUri,"connect err:",connList[i].err)
			continue
		}
		for k := 0; k < ChannelCount; k++ {
			connList[i].NewChannel(false)
		}
	}
	if ConnectTimeOut == 0{
		time.Sleep(999999999 * time.Second)
	}else{
		time.Sleep(time.Duration(ConnectTimeOut) * time.Second)
	}

	resultBackChan <- 1
}
