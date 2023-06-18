// Contains the template framework for the customer Rest API
package model

import (
	"fmt"
	"sync"
	"time"
)

// m is a variable to lock and unlock for builds or concurrency modifications
var m sync.Mutex

// Status represents the Order status of an underlying string type
type Status string

// const for an Order Status type
const (
	Pending  Status = "PENDING"
	Created  Status = "CREATED"
	Payed    Status = "PAYED"
	Shipped  Status = "SHIPPED"
	Received Status = "RECEIVED"
	Canceled Status = "CANCELED"
)

// Customer is a structure representation for a customer's data
type Customer struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Disabled bool    `json:"disabled"`
	Orders   []Order `json:"-"` // remove orders field from serialization
}

// NewCustomer creates a new Customer
func NewCustomer(name string, age int) Customer {
	return Customer{
		Id:     generateId(),
		Name:   name,
		Age:    age,
		Orders: make([]Order, 0),
	}
}

// The Add method adds an order to the customer's order slice. A new order will be added if it does not already exist by order id, otherwise nothing will be added.
func (c *Customer) Add(order Order) {
	if _, contains := c.Contains(order.Id); !contains {
		c.Orders = append(c.Orders, order)
	}
}

// contains method checks if order already exists in customer orders by order id. Return the order if exists
func (c *Customer) Contains(orderId string) (order *Order, exist bool) {
	m.Lock()
	defer m.Unlock()
	for _, o := range c.Orders {
		if o.Id == orderId {
			return &o, true
		}
	}
	return &Order{}, false
}

// Order is a structure representation for the customer order
type Order struct {
	Id      string  `json:"id"`
	Product string  `json:"product"`
	Price   float64 `json:"price"`
	Amount  int     `json:"amount"`
	Status  Status  `json:"status"`
}

// NewOrder creates a new Order
func NewOrder(product string, price float64, amount int) Order {
	return Order{
		Id:      generateId(),
		Product: product,
		Price:   price,
		Amount:  amount,
		Status:  Pending,
	}
}

// generateId is a function to generate value based on machine time in Unix milli precision
// This function used the mutex for concurrency purposes
func generateId() string {
	m.Lock()
	defer m.Unlock()
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}
