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
	_, err := consumer_write.NewBarista()
	k, err := consumer_read.NewKassenpersonal()
	if err != nil {
		fmt.Println(err)
		return
	}

	broker_util, err := broker.NewBrokerUtil("tcp://localhost:1883", "producer", "user", "password")
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

		var choice, drink, size, name, order_id string

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

			broker_util.PublishEvent(broker.TOPIC_ORDER_CREATE, string(payload))

		case "2":
			fmt.Printf("Gebe deine Bestellnummer ein: ")
			fmt.Scan(&order_id)
			id, err := strconv.Atoi(order_id)
			if err != nil {
				fmt.Println("Ungültige Bestellnummer")
			}
			status, err := k.GetStatus(id)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("---------------------------------")
			fmt.Println("Dein Bestellstatus:")
			fmt.Printf("Status: %v\n", status)
			fmt.Println("---------------------------------")
		default:
			fmt.Println("-------------------------------------------------------")
			fmt.Println("Dein Bestellstatus:")
			fmt.Println("\nAktuell gibt es kein neues Status-Update vom Barista.")
			fmt.Println("-------------------------------------------------------")

		}
	}
}
