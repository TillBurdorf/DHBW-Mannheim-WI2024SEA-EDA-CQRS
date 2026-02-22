package consumer_read

import (
	"encoding/json"
	"fmt"

	"CQRS_EDA_TILL/broker"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Kassenpersonal struct {
	broker     *broker.BrokerUtil
	status     map[int]string
	statusChan chan map[int]string
}

func NewKassenpersonal(statusChan chan map[int]string) (*Kassenpersonal, error) {
	brokerUtil, err := broker.NewBrokerUtil("tcp://localhost:1883", "kasse", "user", "password")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}
	k := &Kassenpersonal{
		broker:     brokerUtil,
		status:     make(map[int]string),
		statusChan: statusChan,
	}

	k.broker.SubscribeTopic(broker.TOPIC_ORDER_STATE, k.onOrderCreated)

	return k, nil
}

func (k *Kassenpersonal) onOrderCreated(_ mqtt.Client, msg mqtt.Message) {
	var order_state_event broker.CoffeeOrderStatus
	err := json.Unmarshal(msg.Payload(), &order_state_event)
	if err != nil {
		fmt.Printf("Fehler beim lesen der Bestellung: %v", err)
		return
	}

	k.status[order_state_event.OrderID] = order_state_event.Status

	select {
	case <-k.statusChan:
	default:
	}

	k.statusChan <- k.status
}

func SendToWork(statusChan chan map[int]string) {
	_, err := NewKassenpersonal(statusChan)
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}
