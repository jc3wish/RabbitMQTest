package RabbitMQTest

import (
	"github.com/jc3wish/RabbitMQTest/config"
	"flag"
	"log"
	"time"
	"os"
)

func Start(){
	ConfigFile := flag.String("c", "../etc/config.ini", "Test config file path")
	ConfigKey := flag.String("key", "", "Test config [key]")
	flag.Parse()
	config.LoadConf(*ConfigFile)
	resultBackChan := make(chan int,len(config.MyConf))
	NeenBackCount := 0
	HacBackCount := 0
	if *ConfigKey != ""{
		for k,_ := range config.MyConf{
			if k != *ConfigKey{
				delete(config.MyConf,k)
			}
		}
	}
	TestStartTime := time.Now().Unix()
	for key,m := range config.MyConf{
		if _,ok:=m["Method"];!ok{
			log.Println(key," Method not exsit")
			continue
		}
		switch m["Method"] {
		case "single_send":
			go SingleSend(key,m,resultBackChan)
			NeenBackCount++
			break
		case "single_consume":
			go SingleConsume(key,m,resultBackChan)
			NeenBackCount++
			break
		case "all_write","all_consume","all_write_consume":
			go AllQueueOp(key,m,resultBackChan)
			NeenBackCount++
			break
		default:
			log.Println(key," no Method:",m["Method"])
			break
		}
	}

	if NeenBackCount == 0{
		TestEndTime := time.Now().Unix()
		log.Println("Test Over:",TestEndTime,"Use Time(ms):",TestEndTime-TestStartTime)
		return
	}
	log.Println("NeenBackCount:",NeenBackCount)
	for{
		<- resultBackChan
		HacBackCount++
		if HacBackCount >= NeenBackCount{
			TestEndTime := time.Now().Unix()
			log.Println("Test Over:",TestEndTime,"Use Time(ms):",TestEndTime-TestStartTime)
			os.Exit(0)
		}
	}
}