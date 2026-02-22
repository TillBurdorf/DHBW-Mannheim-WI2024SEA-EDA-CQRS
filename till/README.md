# Getting Started
1. Verwende die Docker Compose Datei, um einen RabbitMQ Broker zu starten
2. Starte die Terminalanwendung mit `go run producer.go`
3. Jetzt kannst du einen Kaffee bestellen und nach der Bestellung den Bestellstatus abrufen

# EDA + CQRS
- die Terminalanwendung sendet ein Event an die MQTT-Topic order/create wenn eine Bestellung erstellt wird
- der Barista hört auf diese Topic und übernimmt die Write Commands
- Bestellungen werden in orders gespeichert, was die Datenbank simulieren soll
- Updates für das Read Model werden auf der Topic order/state veröffentlicht
    - Die veröffentlichung der Updates wird durch eine Go-Routine (Thread) simuliert, welche für jede Bestellung gestartet wird
- das Kassenpersonal ist das Read Model, welches Statusupdates zum Kaffee geben kann und intern eine Map mit dem aktuellen Status hat
-> optimiertes Reading durch die schnelle Map


Falls beim Ausführen Probleme auftreten, enthält der folgende Link eine kurze Demo: https://youtu.be/oUs_92HdHyA
