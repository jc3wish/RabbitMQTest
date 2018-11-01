package RabbitMQTest

import "strconv"

func GetIntDefault(s string,defaultInt int) int{
	intNum,err := strconv.Atoi(s)
	if err!=nil{
		return defaultInt
	}
	return intNum
}


func GetByteBySize(Size int) []byte{
	data := make([]byte,Size)
	for i:=0;i<Size;i++{
		data[i] = 1
		//data=append(data,0)
	}
	return data
}