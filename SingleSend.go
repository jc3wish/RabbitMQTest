package RabbitMQTest

import (
	log "log"
	"time"
	"strings"
	"math/rand"
)

/**
单队列写入操作
*/

func SingleSend(key string,config map[string]string,resultBackChan chan int){
	WriteTimeOut := GetIntDefault(config["WriteTimeOut"],1000)
	DeliveryMode := uint8(GetIntDefault(config["DeliveryMode"],2))

	ConnectCount := GetIntDefault(config["ConnectCount"],1)
	ChannelCount := GetIntDefault(config["ChannelCount"],1)
	ChanneWriteCount := GetIntDefault(config["ChanneWriteCount"],0)
	WaitConfirm := GetIntDefault(config["WaitConfirm"],1)

	var WaitConfirmBool bool

	sizeArr := strings.Split(config["DateSize"], ",")
	sizeLen := len(sizeArr)
	BodyList := make([]*[]byte,sizeLen)
	for k,size:=range sizeArr{
		BodyList[k] = GetByteBySize(GetIntDefault(size,0))
	}

	if WaitConfirm == 1 {
		WaitConfirmBool = true
	}else{
		WaitConfirmBool = false
	}

	ExchangeName := config["ExchangeName"]
	RoutingKey := config["RoutingKey"]
	AmqpUri := config["Uri"]

	OverCount := 0
	NeedWaitCount := ConnectCount*ChannelCount
	ResultChan := make(chan int,NeedWaitCount)

	SendStartTime := time.Now().UnixNano() / 1e6
	log.Println(key,"SingleSend start",SendStartTime)
	for i:=1;i<=ConnectCount;i++ {
		conn := NewConn(AmqpUri)
		if conn.err != nil{
			log.Println(AmqpUri,"connect err:",conn.err)
			continue
		}
		for k := 1; k <= ChannelCount; k++ {
			ch, err := conn.NewChannel(WaitConfirmBool)
			if err != nil {
				ResultChan <- 1
				log.Println(key,"NewChannel err:",k,i,err)
				continue
			}
			ch.SetWriteTimeOut(WriteTimeOut)
			go func(n int,ch *Channel) {
				StartTime:=time.Now().UnixNano() / 1e6
				log.Println(key,"send channel",n,ExchangeName,RoutingKey,"start",StartTime)
				for icount := 0; icount < ChanneWriteCount; icount++ {
					_,err := SendMQ(ch, &ExchangeName, &RoutingKey, &DeliveryMode, BodyList[rand.Intn(sizeLen)])
					if err != nil{
						log.Println(key,"sendMQ",err)
					}
				}
				EndTime := time.Now().UnixNano() / 1e6
				log.Println(key,"send channel",n,"end",EndTime," time(ms):",EndTime-StartTime,ExchangeName,RoutingKey,"sendCount:",ChanneWriteCount)
				ResultChan <- 1
				ch.ch.Close()
			}(k,ch)
		}
	}

	for{
		<-ResultChan
		OverCount++
		if OverCount >= NeedWaitCount{
			break
		}
	}
	SendEndTime := time.Now().UnixNano() / 1e6
	log.Println(key,"SingleSend end",SendEndTime," time(ms):",SendEndTime-SendStartTime)
	resultBackChan <- 1
}