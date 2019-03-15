*RabbitMQTest* 用于RabbitMQ性能测试,可提供对单个队列写入,消费以及对多个队列进行同时读写操作的测试.
可配置连接数,通道数

[![Build Status](https://travis-ci.org/jc3wish/RabbitMQTest.svg?branch=master)](https://travis-ci.org/jc3wish/RabbitMQTest)

#### 编译
```
go build ./RabbitMQTest.go

```

#### 配置参数
```
vim etc/config.ini

[singleSend]
#只写
Method=single_send
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
ExchangeName=amq.direct
RoutingKey=testQueue

#写入等待返回多久
WriteTimeOut=10
#1不持久化  , 2 持久化
DeliveryMode=1
#写入的数据大小,多个用逗号隔开,随机写入,单位byte
DataSize=1024,4096
#连接数
ConnectCount=1
#每个连接开启多个通道
ChannelCount=2
#每个通道总的写入多少条数据
ChanneWriteCount=5000
#开启返回confirm模式
WaitConfirm=1

[singleConsume]
#消费模块
Method=single_consume
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
QueueName=testQueue
#持续多少秒没数据退出,结束通道
ConsumeTimeOut=20
#开启多少个连接,这里不能设置多通道
ConnectCount=1
#消费多少条数据结束通道
ConsumeCount=1000
#是否需要ack数据,这里设置0的话,可以用于测试帐号密码等是否正常
AutoAck=1


[allQueueWrite]
#同时对多个队列操作写操作
Method=all_write

#HttpUri和QueueList 两个参数,QueueList 优先
#QueueList 参数是VHOST:QUEUEU,VHOST2:QUEUEU2 格式,用逗号隔开
#HttpUri 是配置rabbitmq http接口连接,会自动拉取
HttpUri=http://127.0.0.1:15672/api/queues/TESTVHOST
QueueList=TESTVHOST:testQueue,testVhost:testqueue2

#rabbitmq tcp连接ip+端口
AmqpUri=10.4.4.199:5672

#rabbitmq连接帐号
AmqpAdmin=admin

#rabbitmq连接密码
AmqpPwd=admin

#写入操时,则认为false
WriteTimeOut=20

#每个队列写入的连接数量
ConnectCount=1

#每个连接的通道数量
ChannelCount=1

#每个队列的每个通道写入总数
#假如总的有100个队列,ConnectCount = 2,ChannelCount = 3,ChanneWriteCount = 1000 ,则每个队列的写入量 = 2 * 3 * 1000
ChanneWriteCount=1000

#开启cofirm机制
WaitConfirm=1

DataSize=1024

#持久化
DeliveryMode=2

#写入操作的时候 ,采用哪一个交换机,这里routingkey强制采用队列名作为routingkey,并且 交换机和routingkey,队列名的绑定关系 ,得事先绑好
#如果采用默认交换机,则不用填写
#ExchangeName=amq.direct

ContinueCount=80
ContinueCountSleepTime=2


[allQueueConsume]
#多队列消费模式
Method=all_consume
#参考all_write
HttpUri=http://127.0.0.1:15672/api/queues/TESTVHOST
AmqpUri=127.0.0.1:5672
AmqpAdmin=admintest
AmqpPwd=admintest
ConnectCount=1
ConsumeTimeOut=20

#起启多少个队列消费后，sleep [ContinueCountSleepTime] s,0代表不sleep
ContinueCount=80
ContinueCountSleepTime=2


[allQueueWriteAndConsume]
#多队列同时写与消费操作,请参考all_write 和 all_consume 参数
#这里唯一不同的主是 CosumeConnectCount ,WriteConnectCount 分别表示每个队列的消费者数量和每个队列写入开启多少个连接,和 ConnectCount同义,但这里配置 ConnectCount 无效
Method=all_write_consume
HttpUri=http://127.0.0.1:15672/api/queues/TESTVHOST
AmqpUri=127.0.0.1:5672
AmqpAdmin=admintest
AmqpPwd=admintest
WriteTimeOut=20
WriteConnectCount=1
ChannelCount=1
WaitConfirm=1
DataSize=1024
DeliveryMode=2
ChanneWriteCount=1000
ConsumeTimeOut=20
ConsumeCount=1000
CosumeConnectCount=1
AutoAck=1
#ExchangeName=amq.direct
ContinueCount=80
ContinueCountSleepTime=2


[only_connect]
#单连接测试
Method=only_connect
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
ConnectCount=1
ChannelCount=1

#连接完多久退出
ConnectTimeOut=100


[only_declare]
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