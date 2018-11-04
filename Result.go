package RabbitMQTest

type Result struct {
	ConnectSuccess int
	ConnectFail int
	ChannelSuccess int
	ChanneFail int
	WriteSuccess int
	WriteFail int
	CosumeSuccess int
	Type int8
}

func NewResult() *Result{
	return &Result{
		ConnectSuccess:0,
		ConnectFail:0,
		ChannelSuccess:0,
		ChanneFail:0,
		WriteSuccess:0,
		WriteFail:0,
		CosumeSuccess:0,
		Type:0,
	}
}

func (This *Result) setType(Type int8){
	This.Type = Type
}

func (This *Result) addValueConnectSuccess(n int){
	This.ConnectSuccess += n
}
func (This *Result) addValueConnectFail(n int){
	This.ConnectFail += n
}
func (This *Result) addValueChannelSuccess(n int){
	This.ChannelSuccess += n
}
func (This *Result) addValueChanneFail(n int){
	This.ChanneFail += n
}
func (This *Result) addValueWriteSuccess(n int){
	This.WriteSuccess += n
}

func (This *Result) addValueWriteFail(n int){
	This.WriteFail += n
}

func (This *Result) addValueCosumeSuccess(n int){
	This.CosumeSuccess += n
}