package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func GetBrokerStatus(brokers []string) {
	config := sarama.NewConfig()
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal("Error creating Kafka client: ", err)
	}
	defer client.Close()

	for _, broker := range client.Brokers() {
		connected, err := broker.Connected()
		if err != nil {
			log.Printf("Broker %s connection error: %v\n", broker.Addr(), err)
		} else {
			status := "disconnected"
			if connected {
				status = "connected"
			}
			log.Printf("Broker %s is %s\n", broker.Addr(), status)
		}
	}
}
