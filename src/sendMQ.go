package RabbitMQTest

import (
	"github.com/streadway/amqp"
	"fmt"
	"time"
)

func SendMQ(Channel *Channel,exchange *string,routingkey *string,DeliveryMode *uint8,c *[]byte)(bool,error){
	err := Channel.ch.Publish(
		*exchange,     // exchange
		*routingkey, // routing key
		true,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:   *c,
			DeliveryMode:	*DeliveryMode,
			//Expiration:	This.expir,
		})
	if err != nil{
		return false,err
	}
	if Channel.confirmWait != nil{
		select {
		case d := <-Channel.confirmWait:
			if d.DeliveryTag >= 0{
				return true,nil
			}
			return false,fmt.Errorf("unkonw err")
			break
		case <-time.After(time.Duration(Channel.WriteTimeOut) * time.Second):
			return false,fmt.Errorf("server no response")
			break
		}
	}else{
		return true,nil
	}
	return false,fmt.Errorf("unkonw err")
}
