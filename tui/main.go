package tui

import (
	"github.com/rivo/tview"
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
