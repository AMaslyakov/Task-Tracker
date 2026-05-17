package API

import (
	"backend/events"
	"log"
)

var rmqClient *events.RabbitMQClient

// SetRabbitMQClient вызывается из main() для инъекции зависимости
func SetRabbitMQClient(client *events.RabbitMQClient) {
	rmqClient = client
}

// publishEvent — вспомогательная функция для публикации
func publishEvent(evt events.Event) {
	if rmqClient == nil {
		return // dev-mode без RabbitMQ
	}
	if err := rmqClient.Publish(evt); err != nil {
		log.Printf("failed to publish event %s: %v", evt.EventType, err)
	}
}
