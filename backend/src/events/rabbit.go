package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn     *amqp091.Connection
	channel  *amqp091.Channel
	exchange string
	mu       sync.RWMutex
}

func NewRabbitMQClient(ctx context.Context, brokerURL, exchangeName string) (*RabbitMQClient, error) {
	conn, err := amqp091.Dial(brokerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}
	client := &RabbitMQClient{
		conn:     conn,
		channel:  ch,
		exchange: exchangeName,
	}

	go client.monitorConnection(ctx)

	return client, nil
}

func (c *RabbitMQClient) Publish(event Event) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	routingKey := string(event.EventType)

	if err := c.channel.PublishWithContext(
		context.Background(),
		c.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Timestamp:   event.OccurredAt,
			MessageId:   event.EventID,
			Type:        string(event.EventType),
			Body:        body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (c *RabbitMQClient) Consume(ctx context.Context, handler func(Event) error) error {
	q, err := c.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	eventTypes := []EventType{
		EventTaskCreated, EventTaskUpdated, EventTaskStatusChanged,
		EventTaskDeleted, EventTeamCreated, EventTeamSettingsUpdated,
	}
	for _, et := range eventTypes {
		if err := c.channel.QueueBind(
			q.Name,
			string(et),
			c.exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", et, err)
		}
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("consumer channel closed")
			}
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("failed to unmarshal event: %v", err)
				continue
			}
			if err := handler(event); err != nil {
				log.Printf("event handler error: %v", err)
			}
		}
	}
}

func (c *RabbitMQClient) monitorConnection(ctx context.Context) {
	notifyClose := c.conn.NotifyClose(make(chan *amqp091.Error))

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-notifyClose:
			log.Printf("RabbitMQ connection lost: %v, reconnecting...", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *RabbitMQClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
