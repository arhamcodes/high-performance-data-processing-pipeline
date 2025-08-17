package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

// type Address struct {
// 	street   string
// 	city     string
// 	state    string
// 	zip_code string
// 	country  string
// }

// type Customer struct {
// 	id              string
// 	email           string
// 	firstName       string
// 	lastName        string
// 	shippingAddress Address
// 	billingAddress  Address
// }

// type Item struct {
// 	productId string
// 	name      string
// 	price     float64
// 	quantity  int
// 	variant   map[string]string
// }

// type ShippingMethod struct {
// 	id   string
// 	name string
// 	cost float64
// }

// type Payment struct {
// 	method   string
// 	token    string
// 	amount   float64
// 	currency string
// }

// type Order struct {
// 	customer       Customer
// 	items          []Item
// 	shippingMethod ShippingMethod
// 	payment        Payment
// 	orderTotal     float64
// 	taxAmount      float64
// 	discountAmount float64
// 	notes          map[string]string
// }

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Customer struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	ShippingAddress Address `json:"shippingAddress"`
	BillingAddress  Address `json:"billingAddress"`
}

type Item struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	// Variant   map[string]string `json:"variant"`
	Variant string `json:"variant"`
}

type ShippingMethod struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type Payment struct {
	Method   string  `json:"method"`
	Token    string  `json:"token"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Order struct {
	Customer       Customer          `json:"customer"`
	Items          []Item            `json:"items"`
	ShippingMethod ShippingMethod    `json:"shippingMethod"`
	Payment        Payment           `json:"payment"`
	OrderTotal     float64           `json:"orderTotal"`
	TaxAmount      float64           `json:"taxAmount"`
	DiscountAmount float64           `json:"discountAmount"`
	Notes          map[string]string `json:"notes"`
}

type QueueMessage struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Order     Order  `json:"order"`
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"orders",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := context.Background()

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Worker started, waiting for messages...")

	for msg := range msgs {
		var queueMsg QueueMessage
		if err := json.Unmarshal(msg.Body, &queueMsg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		log.Printf("Processing order %s", queueMsg.ID)

		// Process the order (existing logic)
		// ...

		// Store processed result in Redis
		result := map[string]interface{}{
			"id":        queueMsg.ID,
			"timestamp": queueMsg.Timestamp,
			"processed": time.Now().UTC().Format(time.RFC3339),
			"status":    "processed",
			"order":     queueMsg.Order,
		}

		resultJson, _ := json.Marshal(result)
		err = rdb.Set(ctx, "order:"+queueMsg.ID, resultJson, 24*time.Hour).Err()
		if err != nil {
			log.Printf("Error storing in Redis: %v", err)
		}
	}
}
