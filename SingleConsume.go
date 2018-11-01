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

	QueueName := config["QueueName"]
	AmqpUri := config["Uri"]

	OverCount := 0
	NeedWaitCount := ConnectCount
	ResultChan := make(chan int,NeedWaitCount)

	SingleStartTime := time.Now().Unix()
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
			StartTime:=time.Now().Unix()
			log.Println(key,"channel",n,"start",StartTime)
			Cosume(ch,&QueueName)
			EndTime := time.Now().Unix()
			log.Println(key,"channel",n,"end",EndTime," time(ms):",EndTime-StartTime)
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
	SingleEndTime := time.Now().Unix()
	log.Println(key,"SingleConsume end",SingleEndTime," time(ms):",SingleEndTime-SingleStartTime)
	resultBackChan <- 1
}