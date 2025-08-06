package main

import "fmt"

type Address struct {
	street   string
	city     string
	state    string
	zip_code string
	country  string
}

type Customer struct {
	id              string
	email           string
	firstName       string
	lastName        string
	shippingAddress Address
	billingAddress  Address
}

type Item struct {
	productId string
	name      string
	price     float64
	quantity  int
	variant   map[string]string
}

type ShippingMethod struct {
	id   string
	name string
	cost float64
}

type Payment struct {
	method   string
	token    string
	amount   float64
	currency string
}

type Order struct {
	customer       Customer
	items          []Item
	shippingMethod ShippingMethod
	payment        Payment
	orderTotal     float64
	taxAmount      float64
	discountAmount float64
	notes          map[string]string
}

func main() {
	orders := []Order{
		{
			customer: Customer{
				id:        "C001",
				email:     "john.doe@example.com",
				firstName: "John",
				lastName:  "Doe",
				shippingAddress: Address{
					street:   "123 Main St",
					city:     "Metropolis",
					state:    "NY",
					zip_code: "10001",
					country:  "USA",
				},
				billingAddress: Address{
					street:   "123 Main St",
					city:     "Metropolis",
					state:    "NY",
					zip_code: "10001",
					country:  "USA",
				},
			},
			items: []Item{
				{
					productId: "P100",
					name:      "Widget",
					price:     19.99,
					quantity:  2,
					variant:   map[string]string{"color": "red"},
				},
			},
			shippingMethod: ShippingMethod{
				id:   "S1",
				name: "Standard",
				cost: 5.00,
			},
			payment: Payment{
				method:   "Credit Card",
				token:    "tok_abc123",
				amount:   44.98,
				currency: "USD",
			},
			orderTotal:     44.98,
			taxAmount:      3.50,
			discountAmount: 0.00,
			notes:          map[string]string{"gift": "yes"},
		},
	}

	for i, order := range orders {
		println("Processing order", i+1)
		println("Customer:", order.customer.firstName, order.customer.lastName)
		fmt.Printf("Order Total: %.2f\n", order.orderTotal)
	}
}
