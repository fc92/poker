// package groom is used to add or remove poker room
// based on a Welcome screen in terminal mode
package groom

import (
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/rivo/tview"

	"github.com/fc92/poker/internal/player"
)

const (
	newRoomLabel  = "New room name (letters only):"
	openRoomLabel = "open new room"
	inputSize     = 20
)

func DisplayWelcome(rooms []string, serverUrl string) {
	app := tview.NewApplication()

	// Title
	pokerFigure := figure.NewFigure("Poker", "", true)
	title := pokerFigure.String()

	// Form fields
	textView := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).
		SetTextAlign(0).
		SetMaxLines(10).
		SetText(title).
		SetChangedFunc(func() {
			app.Draw()
		})
	nameInput := tview.NewInputField().SetLabel("Enter your name:").SetFieldWidth(inputSize)
	newRoom := tview.NewInputField().SetLabel(newRoomLabel).SetFieldWidth(inputSize)
	rooms = append(rooms, openRoomLabel)
	roomSelection := tview.NewDropDown().
		SetFieldWidth(inputSize).
		SetLabel("Select poker room:")

	// Form
	form := tview.NewForm().
		AddFormItem(nameInput).
		AddFormItem(roomSelection).
		AddButton("Play", func() {
			playerName := strings.TrimSpace(nameInput.GetText())
			if len(playerName) > 0 {
				_, roomSelected := roomSelection.GetCurrentOption()
				if len(roomSelected) > 0 {
					// existing room
					if roomSelected != newRoomLabel {
						startClient(playerName, serverUrl+"/"+roomSelected, app)
					} else {
						// open new room
						newRoomName := strings.TrimSpace(newRoom.GetText())
						if len(newRoomName) > 0 {
							AddRoom(newRoomName)
							startClient(playerName, serverUrl+"/"+newRoomName, app)
						}
					}
				}
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	// Show/Hide new room name in form
	roomSelection.SetOptions(rooms, func(option string, index int) {
		if option == openRoomLabel {
			if form.GetFormItemByLabel(newRoomLabel) == nil {
				form.AddFormItem(newRoom)
			}
			app.SetFocus(newRoom)
		} else {
			if form.GetFormItemByLabel(newRoomLabel) != nil {
				form.RemoveFormItem(form.GetFormItemIndex(newRoomLabel))
			}
			newRoom.SetText("")
		}
	})

	// github link
	githubLink := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).
		SetTextAlign(tview.AlignRight).
		SetMaxLines(10).
		SetText("https://github.com/fc92/poker").
		SetChangedFunc(func() {
			app.Draw()
		})

	// Build screen
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 6, 1, false).
		AddItem(form, 10, 1, true).
		AddItem(githubLink, 1, 1, true)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func startClient(playerName string, serverUrl string, app *tview.Application) {
	app.Stop()
	player.Play(playerName, serverUrl)
}
