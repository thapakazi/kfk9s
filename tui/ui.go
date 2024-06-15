// tui/ui.go
package tui

import (
	"github.com/IBM/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
)

func StartUI() {
	app := tview.NewApplication()

	brokerAddressInput := tview.NewInputField().
		SetLabel("Broker Address").
		SetFieldWidth(20).
		SetText("localhost:9092") //set default address

	// let create a form to take broker address
	form := tview.NewForm().
		AddFormItem(brokerAddressInput).
		AddButton("Submit", func() {
			brokerAddress := brokerAddressInput.GetText()
			displayBrokerStatus(app, brokerAddress)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter Broker Address").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(form, true).Run(); err != nil {
		panic(err)
	}
}

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

func displayTopics(app *tview.Application, brokerAddress string) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	// Set header row for topics
	headers := []string{"TOPIC NAME"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}

	// Function to refresh the table with topics
	refreshTopics := func() {
		// Clear existing rows except the header
		table.Clear().SetCell(0, 0, tview.NewTableCell("TOPIC NAME").
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

		// Fetch topics
		brokers := []string{brokerAddress}
		config := sarama.NewConfig()
		client, err := sarama.NewClient(brokers, config)
		if err != nil {
			log.Printf("Error creating Kafka client: %v\n", err)
			table.SetCell(1, 0, tview.NewTableCell("error").SetAlign(tview.AlignCenter))
		} else {
			defer client.Close()
			topics, err := client.Topics()
			if err != nil {
				log.Printf("Error fetching topics: %v\n", err)
				table.SetCell(1, 0, tview.NewTableCell("error").SetAlign(tview.AlignCenter))
			} else {
				for i, topic := range topics {
					table.SetCell(i+1, 0, tview.NewTableCell(topic).SetAlign(tview.AlignLeft))
				}
			}
		}
	}

	// Call the refresh function initially to populate the table
	refreshTopics()

	// Set up key event handling for refreshing topics and watching for new messages
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'r' {
			refreshTopics()
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'w' {
			row, _ := table.GetSelection()
			if row > 0 {
				topic := table.GetCell(row, 0).Text
				displayMessages(app, brokerAddress, topic)
			}
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			displayBrokerStatus(app, brokerAddress)
			return nil
		}
		return event
	})

	// Display the table
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func displayMessages(app *tview.Application, brokerAddress, topic string) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	// Set header row for messages
	headers := []string{"OFFSET", "PARTITION", "TIMESTAMP", "KEY", "VALUE"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}

	// Function to consume messages from the topic and update the table
	go func() {
		config := sarama.NewConfig()
		config.Consumer.Return.Errors = true
		client, err := sarama.NewConsumer([]string{brokerAddress}, config)
		if err != nil {
			log.Printf("Error creating Kafka consumer: %v\n", err)
			app.Stop()
			return
		}
		defer client.Close()

		partitionConsumer, err := client.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Printf("Error creating partition consumer: %v\n", err)
			app.Stop()
			return
		}
		defer partitionConsumer.Close()

		for message := range partitionConsumer.Messages() {
			app.QueueUpdateDraw(func() {
				row := table.GetRowCount()
				table.SetCell(row, 0, tview.NewTableCell(strconv.FormatInt(message.Offset, 10)).SetAlign(tview.AlignCenter)).
					SetCell(row, 1, tview.NewTableCell(strconv.FormatInt(int64(message.Partition), 10)).SetAlign(tview.AlignCenter)).
					SetCell(row, 2, tview.NewTableCell(message.Timestamp.String()).SetAlign(tview.AlignCenter)).
					SetCell(row, 3, tview.NewTableCell(string(message.Key)).SetAlign(tview.AlignLeft)).
					SetCell(row, 4, tview.NewTableCell(string(message.Value)).SetAlign(tview.AlignLeft))
			})
		}
	}()

	// Set up key event handling to go back to topics view
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			displayTopics(app, brokerAddress)
			return nil
		}
		return event
	})

	// Display the table
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
