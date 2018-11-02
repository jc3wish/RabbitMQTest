package RabbitMQTest

import (
	"log"
	"time"
)

/**
单队消费操作
*/

func SingleConsume(key string,config map[string]string,resultBackChan chan int){
	ConsumeTimeOut := GetIntDefault(config["ConsumeTimeOut"],100)
	ConnectCount := GetIntDefault(config["ConnectCount"],1)
	AutoAck := GetIntDefault(config["AutoAck"],0)

	var ConsumeCount int
	if _,ok:=config["ConsumeCount"];!ok{
		ConsumeCount = 0
	}else{
		ConsumeCount = GetIntDefault(config["ConsumeCount"],0)
	}
	/*
	var AutoAckBool bool
	if AutoAck == 1{
		AutoAckBool = true
	}else{
		AutoAckBool = false
	}
	*/

	QueueName := config["QueueName"]
	AmqpUri := config["Uri"]

	OverCount := 0
	NeedWaitCount := ConnectCount
	ResultChan := make(chan int,NeedWaitCount)

	SingleStartTime := time.Now().UnixNano() / 1e6
	log.Println(key,"SingleConsume start",SingleStartTime)
	for i:=1;i<=ConnectCount;i++ {
		conn := NewConn(AmqpUri)
		if conn.err != nil{
			log.Println(AmqpUri,"connect err:",conn.err)
			continue
		}
		ch, err := conn.NewChannel(false)
		if err != nil {
			ResultChan <- 1
			log.Println(key,"NewChannel err:",i,err)
			continue
		}
		ch.SetConsumeTimeOut(ConsumeTimeOut)
		go func(n int) {
			StartTime:=time.Now().UnixNano() / 1e6
			log.Println(key,"consume channel",n,QueueName,"start",StartTime)
			//Cosume(ch,&QueueName,&ConsumeCount)

			msgs, err := ch.ch.Consume(
				QueueName, // queue
				"",     // consumer
				false,   // auto ack
				false,  // exclusive
				false,  // no local
				false,  // no wait
				nil,    // args
			)
			var HadCosumeCount int = 0
			if err == nil{
				Loop:
				for {
					select {
					case d :=<-msgs:
						if AutoAck == 1{
							d.Ack(true)
						}
						HadCosumeCount++
						if HadCosumeCount >= ConsumeCount && ConsumeCount > 0 {
							ch.ch.Close()
							break Loop
						}
						break
					case <-time.After(time.Duration(ch.ConsumeTimeOut) * time.Second):
						ch.ch.Close()
						break Loop
					}
				}
			}
			EndTime := time.Now().UnixNano() / 1e6
			log.Println(key,"consume channel",n,"end",EndTime," time(ms):",EndTime-StartTime,QueueName,"cosumeCount:",HadCosumeCount)
			ResultChan <- 1
			ch.ch.Close()
			conn.conn.Close()
		}(i)
	}

	for{
		<-ResultChan
		OverCount++
		if OverCount >= NeedWaitCount{
			break
		}
	}
	SingleEndTime := time.Now().UnixNano() / 1e6
	log.Println(key,"SingleConsume end",SingleEndTime," time(ms):",SingleEndTime-SingleStartTime)
	resultBackChan <- 1
}