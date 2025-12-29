package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/prithvirajv06/nimbus-uta/go/notification/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

func NewRabbitMQ(cfg config.RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		cfg.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Declare queue
	_, err = channel.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		r.config.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error) error {
	// Bind queue to exchange
	err := r.channel.QueueBind(
		queueName,
		"#", // routing key pattern
		r.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",
		false, // auto-ack disabled
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				msg.Nack(false, true) // requeue on error
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) DeclareQueue(name string, durable bool) error {
	_, err := r.channel.QueueDeclare(
		name,
		durable,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQ) BindQueue(queueName, routingKey string) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		r.config.Exchange,
		false,
		nil,
	)
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

// Message types
type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// Publisher patterns
type Publisher interface {
	PublishUserCreated(ctx context.Context, userID string, data map[string]interface{}) error
	PublishUserUpdated(ctx context.Context, userID string, data map[string]interface{}) error
	PublishUserDeleted(ctx context.Context, userID string) error
}

type EventPublisher struct {
	rabbitmq *RabbitMQ
}

func NewEventPublisher(rabbitmq *RabbitMQ) *EventPublisher {
	return &EventPublisher{rabbitmq: rabbitmq}
}

func (p *EventPublisher) PublishUserCreated(ctx context.Context, userID string, data map[string]interface{}) error {
	msg := Message{
		Type:    "user.created",
		Payload: data,
	}
	return p.rabbitmq.Publish(ctx, "user.created", msg)
}

func (p *EventPublisher) PublishUserUpdated(ctx context.Context, userID string, data map[string]interface{}) error {
	msg := Message{
		Type:    "user.updated",
		Payload: data,
	}
	return p.rabbitmq.Publish(ctx, "user.updated", msg)
}

func (p *EventPublisher) PublishUserDeleted(ctx context.Context, userID string) error {
	msg := Message{
		Type: "user.deleted",
		Payload: map[string]interface{}{
			"user_id": userID,
		},
	}
	return p.rabbitmq.Publish(ctx, "user.deleted", msg)
}
