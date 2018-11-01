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
Method=single_send
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
ExchangeName=amq.direct
RoutingKey=testQueue
WriteTimeOut=10
DeliveryMode=2
DateSize=1024
ConnectCount=1
ChannelCount=1
ChanneWriteCount=1000
WaitConfirm=1

[singleConsume]
Method=single_consume
Uri=amqp://admintest:admintest@127.0.0.1:5672/testvhost
QueueName=testQueue
ConsumeTimeOut=20
ConnectCount=1

[allQueueWrite]
Method=all_write
HttpUri=http://127.0.0.1:15672/api/queues/TESTVHOST
AmqpUri=10.4.4.199:5672
AmqpAdmin=admin
AmqpPwd=admin
WriteTimeOut=20
ConnectCount=1
ChannelCount=1
WaitConfirm=1
DateSize=1024
DeliveryMode=2
ChanneWriteCount=1000
ConsumeTimeOut=20

[allQueueConsume]
Method=all_consume
HttpUri=http://127.0.0.1:15672/api/queues/TESTVHOST
AmqpUri=127.0.0.1:5672
AmqpAdmin=admintest
AmqpPwd=admintest
ConnectCount=1
ConsumeTimeOut=20

[allQueueWriteAndConsume]
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
CosumeConnectCount=1

```
#### 运行
```
./StartTest -c ../etc/config.ini -key singleSend

```