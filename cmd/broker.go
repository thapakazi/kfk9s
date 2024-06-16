package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thapakazi/kfk9s/kafka"
)

var brokerCmd = &cobra.Command{
	Use:   "brokers",
	Short: "List Kafka brokers",
	Run: func(cmd *cobra.Command, args []string) {
		brokers := []string{"kafka-production-broker1:9092"} // Replace with your brokers
		kafka.GetBrokerStatus(brokers)
	},
}

func init() {
	rootCmd.AddCommand(brokerCmd)
}
