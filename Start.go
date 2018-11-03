package RabbitMQTest

import (
	"github.com/jc3wish/RabbitMQTest/config"
	"flag"
	"log"
	"time"
	"os"
	"strings"
)

func Start(){
	ConfigFile := flag.String("c", "../etc/config.ini", "Test config file path")
	ConfigKey := flag.String("key", "", "Test config [key]")
	flag.Parse()
	config.LoadConf(*ConfigFile)
	resultBackChan := make(chan int,len(config.MyConf))
	NeenBackCount := 0
	HacBackCount := 0
	var myConf map[string]map[string]string
	if *ConfigKey != ""{
		myConf = make(map[string]map[string]string,0)
		for _,v := range strings.Split(*ConfigKey,","){
			if _,ok:=config.MyConf[v];ok{
				myConf[v] = config.MyConf[v];
			}else{
				log.Println("key:",v," not in config")
			}
		}
	}else{
		myConf=config.MyConf
	}
	TestStartTime := time.Now().UnixNano() / 1e6
	for key,m := range myConf{
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
		case "only_connect":
			go OnlyConnect(key,m,resultBackChan)
			NeenBackCount++
			break
		default:
			log.Println(key," no Method:",m["Method"])
			break
		}
	}

	if NeenBackCount == 0{
		TestEndTime := time.Now().UnixNano() / 1e6
		log.Println("Test Over:",TestEndTime,"Use Time(ms):",TestEndTime-TestStartTime)
		return
	}
	log.Println("NeenBackCount:",NeenBackCount)
	for{
		<- resultBackChan
		HacBackCount++
		if HacBackCount >= NeenBackCount{
			TestEndTime := time.Now().UnixNano() / 1e6
			log.Println("Test Over:",TestEndTime,"Use Time(ms):",TestEndTime-TestStartTime)
			os.Exit(0)
		}
	}
}