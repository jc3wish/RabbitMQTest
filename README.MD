*RabbitMQTest* 用于RabbitMQ性能测试,可提供对单个队列写入,消费以及对多个队列进行同时读写操作的测试.
可配置连接数,通道数

[![Build Status](https://travis-ci.org/jc3wish/RabbitMQTest.svg?branch=master)](https://travis-ci.org/jc3wish/RabbitMQTest)

#### 编译
```
go build ./RabbitMQTest.go

```

#### 配置参数

[**etc/config.ini**](https://github.com/jc3wish/RabbitMQTest/blob/master/etc/config.ini)


#### 案例 1

> *随机生成100个队列及绑定交换机*

##### 配置信息
```
[declareQueue]
#单纯创建队列
Method=only_declare
Uri=amqp://guest:guest@127.0.0.1:35673/mytest

#QueueList参数存在的情况下 QueuePrefix,QueueDurable,QueueAutoDelete,QueueCount参数无效
#QueueList = name:durable,auto_delete
#QueueList=TestQueue1:true:false,TestQueue2:true:false

#自动生成的队列的前缀
QueuePrefix=TestQueue
QueueDurable=true
QueueAutoDelete=false

#生成多少个队列
QueueCount=100

#交换机配置信息
ExchangeName=amq.direct
ExchangeType=direct
ExchangeDurable=true
ExchangeAutoDelete=false

#routingkey默认为队列名,如果填写RoutingKey,则统一绑定这个routingkey
RoutingKey=

```
##### 运行

```
./RabbitMQTest -c ./etc/config.ini -key declareQueue

++++++++++++++++++++++++++++++
+                            +
+        RabbitMQTest        +
+                            +      Version:RabbitMQTest-v1.0.1-release
++++++++++++++++++++++++++++++           By:jc3wish


2019/03/15 18:57:38 Test Start, Time: 1552647458450
2019/03/15 18:57:38 only_declare Declare Count: 100 endTime: 1552647458460  had use time(ms): 10

2019/03/15 18:57:38 Test Over: 1552647458460 Use Time(ms): 10
ConnectSuccess: 0
ConnectFail: 0
ChannelSuccess: 0
ChanneFail: 0
WriteSuccess: 0
WriteFail: 0
CosumeSuccess: 0
Write QPS: 0
Consume QPS: 0

```


#### 案例 2

- 往 vhost:mytest 下所有队列写入数据
- 每个队列一个连接
- 每个连接2个通道
- 每个通道往一个队列中写入1000条数据。
- 每条数据随机大小 1kb,2kb

##### 配置信息
```
[allQueueWrite]
#同时对多个队列操作写操作
Method=all_write

#HttpUri和QueueList 两个参数,QueueList 优先
#QueueList 参数是VHOST:QUEUEU,VHOST2:QUEUEU2 格式,用逗号隔开
#HttpUri 是配置rabbitmq http接口连接,会自动拉取
HttpUri=http://127.0.0.1:15674/api/queues/mytest
#QueueList=TESTVHOST:testQueue,testVhost:testqueue2

#rabbitmq tcp连接ip+端口
AmqpUri=127.0.0.1:35673

#rabbitmq连接帐号
AmqpAdmin=guest

#rabbitmq连接密码
AmqpPwd=guest

#写入操时,则认为false
WriteTimeOut=20

#每个队列写入的连接数量
ConnectCount=1

#每个连接的通道数量
ChannelCount=2

#每个队列的每个通道写入总数
#假如总的有100个队列,ConnectCount = 1,ChannelCount = 2,ChanneWriteCount = 1000 ,则每个队列的写入量 = 1 * 2 * 10
ChanneWriteCount=10

#开启cofirm机制
WaitConfirm=1

#数据大小
DataSize=1024,2048

#持久化
DeliveryMode=1

#写入操作的时候 ,采用哪一个交换机,这里routingkey强制采用队列名作为routingkey,并且 交换机和routingkey,队列名的绑定关系 ,得事先绑好
#如果采用默认交换机,则不用填写
#ExchangeName=amq.direct

ContinueCount=100
ContinueCountSleepTime=2
```

##### 运行

```
./RabbitMQTest -c ./etc/config.ini -key allQueueWrite

++++++++++++++++++++++++++++++
+                            +
+        RabbitMQTest        +
+                            +      Version:RabbitMQTest-v1.0.1-release
++++++++++++++++++++++++++++++           By:jc3wish


2019/03/15 19:03:16 Test Start, Time: 1552647796908
2019/03/15 19:03:16 allQueueWrite AllQueueOp start 1552647796959

2019/03/15 19:03:33 Test Over: 1552647813853 Use Time(ms): 16945
ConnectSuccess: 400
ConnectFail: 0
ChannelSuccess: 800
ChanneFail: 0
WriteSuccess: 8000
WriteFail: 0
CosumeSuccess: 0
Write QPS: 472.115668338743
Consume QPS: 0

```

#### 备注

假如要同时运行多个 key 用逗号隔，例如：

```
./RabbitMQTest -c ./etc/config.ini -key singleSend,singleConsume

```