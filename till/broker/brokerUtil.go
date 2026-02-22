package broker

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewBrokerUtil(url string, clientID string, user string, password string) (*BrokerUtil, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID(clientID)
	opts.SetUsername(user)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &BrokerUtil{mqttClient: client}, nil
}

func (b *BrokerUtil) SubscribeTopic(topic string, message_callback mqtt.MessageHandler) {
	token := b.mqttClient.Subscribe(topic, 1, message_callback)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Fehler! Topic %s nicht subscribed! Fehler: %s", topic, token.Error())
	}
}

func (b *BrokerUtil) PublishEvent(topic string, msg string) {
	token := b.mqttClient.Publish(topic, 1, false, msg)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Fehler beim veröffentlichen: %s", token.Error())
	}
}
