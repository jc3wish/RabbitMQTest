package RabbitMQTest

import (
	log "log"
	"strings"
	"math/rand"
)

/**
单队列写入操作
*/

func SingleSend(key string,config map[string]string,resultDataChan chan *Result){
	ResultData := NewResult()
	ResultData.Type = 1
	WriteTimeOut := GetIntDefault(config["WriteTimeOut"],1000)
	DeliveryMode := uint8(GetIntDefault(config["DeliveryMode"],2))

	ConnectCount := GetIntDefault(config["ConnectCount"],1)
	ChannelCount := GetIntDefault(config["ChannelCount"],1)
	ChanneWriteCount := GetIntDefault(config["ChanneWriteCount"],0)
	WaitConfirm := GetIntDefault(config["WaitConfirm"],1)

	var WaitConfirmBool bool

	sizeArr := strings.Split(config["DataSize"], ",")
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
	NeedWaitCount := 0
	ResultChan := make(chan int,ConnectCount*ChannelCount)

	//SendStartTime := time.Now().UnixNano() / 1e6
	//log.Println(key,"SingleSend start",SendStartTime)
	for i:=1;i<=ConnectCount;i++ {
		conn := NewConn(AmqpUri)
		if conn.err != nil{
			log.Println(AmqpUri,"connect err:",conn.err)
			ResultData.ConnectFail++
			continue
		}
		ResultData.ConnectSuccess++
		for k := 1; k <= ChannelCount; k++ {
			ch, err := conn.NewChannel(WaitConfirmBool)
			if err != nil {
				log.Println(key,"NewChannel err:",k,i,err)
				ResultData.ChanneFail++
				continue
			}
			ResultData.ChannelSuccess++
			ch.SetWriteTimeOut(WriteTimeOut)
			go func(n int,ch *Channel) {
				FailCount:=0;
				//StartTime:=time.Now().UnixNano() / 1e6
				//log.Println(key,"send channel",n,ExchangeName,RoutingKey,"start",StartTime)
				for icount := 0; icount < ChanneWriteCount; icount++ {
					_,err := SendMQ(ch, &ExchangeName, &RoutingKey, &DeliveryMode, BodyList[rand.Intn(sizeLen)])
					if err != nil{
						log.Println(key,"sendMQ",err)
						FailCount++
						continue
					}
				}
				//EndTime := time.Now().UnixNano() / 1e6
				//log.Println(key,"send channel",n,"end",EndTime," time(ms):",EndTime-StartTime,ExchangeName,RoutingKey,"sendCount:",ChanneWriteCount)
				ResultChan <- FailCount
				ch.ch.Close()
			}(k,ch)
			NeedWaitCount++
		}
	}

	if NeedWaitCount > 0 {
		for{
			FailCount := <-ResultChan
			ResultData.WriteFail += FailCount
			ResultData.WriteSuccess += ChanneWriteCount - FailCount
			OverCount++
			if OverCount >= NeedWaitCount {
				break
			}
		}
		/*
		loop:
		for {
			select {
			case FailCount := <-ResultChan:
				ResultData.WriteFail += FailCount
				ResultData.WriteSuccess += ChanneWriteCount - FailCount
				OverCount++
				if OverCount >= NeedWaitCount {
					break loop
				}
			case <-time.After(100 * time.Second):
				log.Println(key, "single Send time after,had use time:", time.Now().UnixNano()/1e6-SendStartTime)
				fmt.Println("ConnectSuccess:", ResultData.ConnectSuccess)
				fmt.Println("ConnectFail:", ResultData.ConnectFail)
				fmt.Println("ChannelSuccess:", ResultData.ChannelSuccess)
				fmt.Println("ChanneFail:", ResultData.ChanneFail)
				fmt.Println("WriteSuccess:", ResultData.WriteSuccess)
				fmt.Println("WriteFail:", ResultData.WriteFail)
				fmt.Println("CosumeSuccess:", ResultData.CosumeSuccess)
				break
			}
		}
		*/
	}
	//SendEndTime := time.Now().UnixNano() / 1e6
	//log.Println(key,"SingleSend end",SendEndTime," time(ms):",SendEndTime-SendStartTime)
	resultDataChan <- ResultData
}