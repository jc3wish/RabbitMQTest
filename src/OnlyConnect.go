package RabbitMQTest

import (
	"log"
	"time"
)

func OnlyConnect(key string,config map[string]string,resultDataChan chan *Result){
	ResultData := NewResult()
	AmqpUri := config["Uri"]
	ConnectCount := GetIntDefault(config["ConnectCount"],1)
	ChannelCount := GetIntDefault(config["ChannelCount"],0)
	ConnectTimeOut := GetIntDefault(config["ConnectTimeOut"],100)
	connList := make([]*Conn,ConnectCount)
	for i:=0;i<ConnectCount;i++{
		connList[i] = NewConn(AmqpUri)
		if connList[i].err != nil{
			ResultData.ConnectFail++
			log.Println(AmqpUri,"connect err:",connList[i].err)
			continue
		}
		ResultData.ConnectSuccess++
		for k := 0; k < ChannelCount; k++ {
			_,err:=connList[i].NewChannel(false)
			if err != nil{
				ResultData.ChanneFail++
				continue
			}
			ResultData.ChannelSuccess++
		}
	}
	if ConnectTimeOut == 0{
		time.Sleep(999999999 * time.Second)
	}else{
		time.Sleep(time.Duration(ConnectTimeOut) * time.Second)
	}

	resultDataChan <- ResultData
}
