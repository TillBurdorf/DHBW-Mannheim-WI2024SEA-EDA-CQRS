package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"CQRS_EDA_TILL/broker"
	"CQRS_EDA_TILL/consumer_read"
	"CQRS_EDA_TILL/consumer_write"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

func main() {
	go consumer_write.SendToWork()

	statusChan := make(chan map[int]string, 1)
	go consumer_read.SendToWork(statusChan)

	brokerUtil, err := broker.NewBrokerUtil("tcp://localhost:1883", "producer", "user", "password")
	if err != nil {
		fmt.Println(err)
		return
	}

	id_counter := 0
	for {

		menu := map[string]string{
			"1": "Cappuchinno",
			"2": "Flat White",
			"3": "Espresso",
			"4": "Americano",
		}

		var choice, drink, size, name, orderID string

		fmt.Printf("Hallo! Wähle bitte eine Option:\n\n(1) Neue Bestellung\n(2) Bestellstatus abrufen\n(0) Beenden\n\nEingabe: ")
		fmt.Scan(&choice)
		switch choice {
		case "0":
			fmt.Println("Tschüss! Bis zum nächsten mal...")
			return
		case "1":
			fmt.Printf("\n--- Neue Bestellung ---\n\n")
			fmt.Printf("(1) Cappuchino\n(2) Flat White\n(3) Espresso\n(4) Americano\n\n")
			fmt.Print("Eingabe (Zahl): ")
			fmt.Scan(&drink)
			drinkName, ok := menu[drink]
			if !ok {
				fmt.Println("Ungültige Wahl!")
			}
			fmt.Print("Größe (S/M/L): ")
			fmt.Scan(&size)
			fmt.Print("Dein Name: ")
			fmt.Scan(&name)
			fmt.Println("---------------------------------")
			fmt.Println("Vielen Dank für deine Bestellung!")
			fmt.Printf("Bestellnummer: %v\n", id_counter)
			fmt.Println("---------------------------------")
			fmt.Printf("\n\n\n")
			orderCreationEvent := broker.CoffeeOrderCreated{
				OrderID: id_counter,
				Drink:   drinkName,
				Size:    size,
				Name:    name,
			}

			id_counter++

			payload, err := json.Marshal(orderCreationEvent)
			if err != nil {
				fmt.Println(err)
				return
			}

			brokerUtil.PublishEvent(broker.TOPIC_ORDER_CREATE, string(payload))

		case "2":
			fmt.Printf("Gebe deine Bestellnummer ein: ")
			fmt.Scan(&orderID)
			id, err := strconv.Atoi(orderID)
			if err != nil {
				fmt.Println("Ungülrige Bestellnummer")
			}
			select {
			case status := <-statusChan:
				currentStatus := status[id]
				statusColor := ColorReset

				switch currentStatus {
				case "Abholbereit":
					statusColor = ColorGreen
				case "In Zubereitung", "Eingegangen":
					statusColor = ColorYellow
				default:
					statusColor = ColorCyan
				}
				fmt.Println("---------------------------------")
				fmt.Println("Dein Bestellstatus:")
				fmt.Printf("Status: %s%v%s\n", statusColor, status[id], ColorReset)
				fmt.Println("---------------------------------")
			default:
				fmt.Println("-------------------------------------------------------")
				fmt.Println("Dein Bestellstatus:")
				fmt.Println("\nAktuell gibt es kein neues Status-Update vom Barista.")
				fmt.Println("-------------------------------------------------------")

			}
		}
	}
}
