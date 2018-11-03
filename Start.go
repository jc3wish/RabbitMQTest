package RabbitMQTest

import (
	"github.com/jc3wish/RabbitMQTest/config"
	"flag"
	"log"
	"time"
	"os"
	"strings"
	"fmt"
)

const VERSION  = "RabbitMQTest-v1.0.0-beta.02"

var (
	v bool
)

var welcome string = `
++++++++++++++++++++++++++++++   
+                            +
+        RabbitMQTest        +
+                            +      Version:{$version}
++++++++++++++++++++++++++++++           By:jc3wish

`

func Start(){

	ConfigFile := flag.String("c", "../etc/config.ini", "Test config file path")
	ConfigKey := flag.String("key", "", "Test config [key]")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.Parse()
	if v {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	welcome = strings.Replace(welcome,"{$version}",VERSION,-1)
	fmt.Println(welcome)

	config.LoadConf(*ConfigFile)

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

	resultDataChan := make(chan *Result,len(myConf))
	TestStartTime := time.Now().UnixNano() / 1e6
	log.Println("Test Start, Time:",TestStartTime)
	for key,m := range myConf{
		if _,ok:=m["Method"];!ok{
			log.Println(key," Method not exsit")
			continue
		}
		switch m["Method"] {
		case "single_send":
			go SingleSend(key,m,resultDataChan)
			NeenBackCount++
			break
		case "single_consume":
			go SingleConsume(key,m,resultDataChan)
			NeenBackCount++
			break
		case "all_write","all_consume","all_write_consume":
			go AllQueueOp(key,m,resultDataChan)
			NeenBackCount++
			break
		case "only_connect":
			go OnlyConnect(key,m,resultDataChan)
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

	ResultData := NewResult()
	for{
		data := <- resultDataChan
		ResultData.ConnectSuccess += data.ConnectSuccess
		ResultData.ConnectFail += data.ConnectFail
		ResultData.ChannelSuccess += data.ChannelSuccess
		ResultData.ChanneFail += data.ChanneFail
		ResultData.WriteSuccess += data.WriteSuccess
		ResultData.WriteFail += data.WriteFail
		ResultData.CosumeSuccess += data.CosumeSuccess
		HacBackCount++
		if HacBackCount >= NeenBackCount{
			TestEndTime := time.Now().UnixNano() / 1e6
			fmt.Println(" ")
			UseTime := float64(TestEndTime-TestStartTime)
			log.Println("Test Over:",TestEndTime,"Use Time(ms):",UseTime)
			fmt.Println("ConnectSuccess:",ResultData.ConnectSuccess)
			fmt.Println("ConnectFail:",ResultData.ConnectFail)
			fmt.Println("ChannelSuccess:",ResultData.ChannelSuccess)
			fmt.Println("ChanneFail:",ResultData.ChanneFail)
			fmt.Println("WriteSuccess:",ResultData.WriteSuccess)
			fmt.Println("WriteFail:",ResultData.WriteFail)
			fmt.Println("CosumeSuccess:",ResultData.CosumeSuccess)
			fmt.Println("Write QPS:",float64(ResultData.WriteSuccess)/UseTime*1000)
			fmt.Println("Consume QPS:",float64(ResultData.CosumeSuccess)/UseTime*1000)
			os.Exit(0)
		}
	}
}