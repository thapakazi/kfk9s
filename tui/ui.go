// tui/ui.go
package tui

import (
    "github.com/rivo/tview"
	"log"
	"github.com/IBM/sarama"
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
		AddButton("Submit", func(){
			brokerAddress := brokerAddressInput.GetText()
			displayBrokerStatus(app,brokerAddress)
		}).
		AddButton("Quit", func(){
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter Broker Address").SetTitleAlign(tview.AlignLeft)

    if err := app.SetRoot(form, true).Run(); err != nil {
        panic(err)
    }
}

func displayBrokerStatus(app *tview.Application, brokerAddress string){
	
    table := tview.NewTable().
	SetFixed(1,1).
		SetSelectable(true,false)

	headers := []string{"NAME","STATUS"}
	for i, header := range headers{
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		table.SetCell(0,i,cell)
	}

    //table.SetCell(1, 0, tview.NewTableCell(brokerAddress)).
    //    SetCell(1, 1, tview.NewTableCell("connected"))

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
	
	// Display the table
    if err := app.SetRoot(table, true).Run(); err != nil {
        panic(err)
    }
}
