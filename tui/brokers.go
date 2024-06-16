package tui

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/rivo/tview"
)

func displayBrokerStatus(app *tview.Application, brokerAddress string) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	// Set header row
	headers := []string{"NAME", "STATUS"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}

	// Fetch broker status
	brokers := []string{brokerAddress}
	config := sarama.NewConfig()
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Printf("Error creating Kafka client: %v\n", err)
		table.SetCell(1, 0, tview.NewTableCell(brokerAddress).SetAlign(tview.AlignLeft)).
			SetCell(1, 1, tview.NewTableCell("error").SetAlign(tview.AlignCenter))
	} else {
		defer client.Close()
		for i, broker := range client.Brokers() {
			err := broker.Open(config)
			if err != nil {
				log.Printf("Error opening connection to broker %s: %v\n", broker.Addr(), err)
				table.SetCell(i+1, 0, tview.NewTableCell(broker.Addr()).SetAlign(tview.AlignLeft)).
					SetCell(i+1, 1, tview.NewTableCell("error").SetAlign(tview.AlignCenter))
				continue
			}

			connected, err := broker.Connected()
			if err != nil {
				log.Printf("Broker %s connection error: %v\n", broker.Addr(), err)
				table.SetCell(i+1, 0, tview.NewTableCell(broker.Addr()).SetAlign(tview.AlignLeft)).
					SetCell(i+1, 1, tview.NewTableCell("error").SetAlign(tview.AlignCenter))
			} else {
				status := "disconnected"
				if connected {
					status = "connected"
				}
				table.SetCell(i+1, 0, tview.NewTableCell(broker.Addr()).SetAlign(tview.AlignLeft)).
					SetCell(i+1, 1, tview.NewTableCell(status).SetAlign(tview.AlignCenter))
			}
		}
	}

	// Set up selection handling
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 {
			brokerAddr := table.GetCell(row, 0).Text
			displayTopics(app, brokerAddr)
		}
	})

	// Display the table
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
