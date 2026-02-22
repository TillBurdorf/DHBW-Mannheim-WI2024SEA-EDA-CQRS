package consumer_write

import (
	"encoding/json"
	"fmt"
	"time"

	"CQRS_EDA_TILL/broker"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Barista struct {
	broker *broker.BrokerUtil
	orders []broker.CoffeeOrderCreated
}

func NewBarista() (*Barista, error) {
	util, err := broker.NewBrokerUtil("tcp://localhost:1883", "barista", "user", "password")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	b := &Barista{
		broker: util,
		orders: []broker.CoffeeOrderCreated{},
	}

	b.broker.SubscribeTopic(broker.TOPIC_ORDER_CREATE, b.onOrderCreated)

	return b, nil
}

func (b *Barista) onOrderCreated(_ mqtt.Client, msg mqtt.Message) {
	var coffe_order_created_event broker.CoffeeOrderCreated
	err := json.Unmarshal(msg.Payload(), &coffe_order_created_event)
	if err != nil {
		fmt.Printf("Fehler beim lesen der Bestellung: %v", err)
		return
	}

	b.orders = append(b.orders, coffe_order_created_event)
	go b.giveUpdates(coffe_order_created_event)
}

func (b *Barista) giveUpdates(coffe_order_created_event broker.CoffeeOrderCreated) {
	// Simuliere Bestellungseingang
	status := broker.CoffeeOrderStatus{
		OrderID: coffe_order_created_event.OrderID,
		Name:    coffe_order_created_event.Name,
		Status:  "Eingegangen",
	}
	b.publishStatus(status)
	time.Sleep(time.Second * 10)

	// Simuliere Bestellungsvorbereitung
	status.Status = "Wird zubereitet"
	b.publishStatus(status)
	time.Sleep(time.Second * 10)

	// Simuliere Bestellungsabschluss
	status.Status = "Abholbereit"
	b.publishStatus(status)
	time.Sleep(time.Second * 10)
}

func (b *Barista) publishStatus(status broker.CoffeeOrderStatus) {
	payload, err := json.Marshal(status)
	if err != nil {
		fmt.Println(err)
		return
	}
	b.broker.PublishEvent(broker.TOPIC_ORDER_STATE, string(payload))
}
