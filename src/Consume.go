package RabbitMQTest

import (
	"fmt"
	"time"
)


func Cosume(Channel *Channel,QueueName *string,ConsumeCount *int) error{
	msgs, err := Channel.ch.Consume(
		*QueueName, // queue
		"",     // consumer
		false,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil{
		return err
	}
	var HadCosumeCount int = 0
	for{
		select {
		case	d := <-msgs:
			d.Ack(false)
			HadCosumeCount++
			if HadCosumeCount >= *ConsumeCount && *ConsumeCount>0{
				Channel.ch.Close()
				return nil
			}
			break
		case <-time.After(time.Duration(Channel.ConsumeTimeOut) * time.Second):
			Channel.ch.Close()
			return fmt.Errorf("ConsumeTimeOut ")
		}
	}
	return nil
}

