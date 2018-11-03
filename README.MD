*RabbitMQTest* 用于RabbitMQ性能测试,可提供对单个队列写入,消费以及对多个队列进行同时读写操作的测试.
可配置连接数,通道数

#### 编译
```
go build ./Start/StartTest.go

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
DeliveryMode=2
#写入的数据大小,多个用逗号隔开,随机写入,单位byte
DateSize=1024,4096
#连接数
ConnectCount=1
#每个连接开启多个通道
ChannelCount=1
#每个通道总的写入多少条数据
ChanneWriteCount=1000
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

DateSize=1024

#持久化
DeliveryMode=2

#写入操作的时候 ,采用哪一个交换机,这里routingkey强制采用队列名作为routingkey,并且 交换机和routingkey,队列名的绑定关系 ,得事先绑好
ExchangeName=amq.direct


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
DateSize=1024
DeliveryMode=2
ChanneWriteCount=1000
ConsumeTimeOut=20
ConsumeCount=1000
CosumeConnectCount=1
AutoAck=1
ExchangeName=amq.direct


[only_connect]
#单连接测试
Method=only_connect
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
ConnectCount=1
ChannelCount=1

#连接完多久退出
ConnectTimeOut=100

```
#### 运行
```
./StartTest -c ../etc/config.ini -key singleSend,singleConsume

```