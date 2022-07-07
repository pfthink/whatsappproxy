package rabbitmq

import "sync"

var (
	imRabbitMq *RabbitMQ
	syOnce     sync.Once
)

func InitBossRabbitMq() *RabbitMQ {
	if imRabbitMq == nil {
		syOnce.Do(func() {
			queueExchange := &QueueExchange{
				"im.sendMessage.whatsapp",
				"im.sendMessage.whatsapp",
				"im.sendMessage.whatsapp",
				"topic",
			}
			imRabbitMq = New(queueExchange)
		})
	}
	return imRabbitMq
}

func SendBossImMsg(msg []byte) {
	t := &Message{
		msg,
	}
	go imRabbitMq.listenProducer(t)
}
