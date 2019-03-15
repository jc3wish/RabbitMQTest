package RabbitMQTest

import (
	"strconv"
	"sync"
)

var sizeMap map[int]*[]byte
var SizeLock sync.Mutex

func init(){
	sizeMap = make(map[int]*[]byte,0)
}

func GetIntDefault(s string,defaultInt int) int{
	intNum,err := strconv.Atoi(s)
	if err!=nil{
		return defaultInt
	}
	return intNum
}

func GetByteBySize(Size int) *[]byte{
	SizeLock.Lock()
	defer SizeLock.Unlock()
	if _,ok:=sizeMap[Size];ok{
		return sizeMap[Size]
	}
	data := make([]byte,Size)
	for i:=0;i<Size;i++{
		data[i] = 1
	}
	sizeMap[Size] = &data;
	return &data
}