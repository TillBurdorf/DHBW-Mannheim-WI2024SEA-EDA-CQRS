package broker

import mqtt "github.com/eclipse/paho.mqtt.golang"

type CoffeeOrderCreated struct {
	OrderID int    `json:"order_id"`
	Name    string `json:"name"`
	Drink   string `json:"drink"`
	Size    string `json:"size"`
}

type CoffeeOrderStatus struct {
	OrderID int    `json:"order_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type BrokerUtil struct {
	mqttClient mqtt.Client
}

const (
	TOPIC_ORDER_CREATE string = "order/create"
	TOPIC_ORDER_STATE  string = "order/status"
)
