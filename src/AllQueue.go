package RabbitMQTest

import (
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type queueInfo struct {
	Vhost string `json:"vhost"`
	Queue string `json:"name"`
}

func GetQueuesByUrl(url string,user string,pwd string) *[]queueInfo{
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization","Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pwd)))
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var qList []queueInfo
	qList = make([]queueInfo,0)
	json.Unmarshal(body,&qList)
	return &qList
}

func GetQueuesByConfig(QeueuList string) *[]queueInfo{
	s := strings.Split(QeueuList,",")
	qList := make([]queueInfo,0)
	for _,v := range s{
		q := strings.Split(v,":")
		qList = append(qList, queueInfo{Vhost:q[0],Queue:q[1]});
	}
	return &qList
}
