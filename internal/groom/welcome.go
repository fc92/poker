// package groom is used to add or remove poker room
// based on a Welcome screen in terminal mode
package groom

import (
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	newRoomLabel  = "New room name (letters only):"
	openRoomLabel = "open new room"
	urlLabel      = "Use this url to join the poker room: "
	inputSize     = 20
)

// Terminal welcome page to choose player name and poker room
// serverUrl designate the url used by the browser to reach the poker room
func DisplayWelcome(serverUrl string) {
	rooms, err := RoomDeployed()
	if err != nil {
		log.Error().Msg("unable to get list of rooms deployed...")
	} else {
		log.Debug().Msgf("Found rooms: %v", rooms)
	}
	rooms = append(rooms, openRoomLabel)
	app := tview.NewApplication()
	roomUrl := ""
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Title
	pokerFigure := figure.NewFigure("Team Poker", "", true)
	title := pokerFigure.String()

	// github link
	githubLink := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).
		SetTextAlign(tview.AlignRight).
		SetMaxLines(1).
		SetText("https://github.com/fc92/poker").
		SetChangedFunc(func() {
			app.Draw()
		})

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
	displayUrl := tview.NewTextView().SetLabel(urlLabel)
	roomSelection := tview.NewDropDown().
		SetFieldWidth(inputSize).
		SetLabel("Select poker room:")

	// Form
	form := tview.NewForm().
		AddFormItem(nameInput).
		AddFormItem(roomSelection).
		AddButton("Get room url", func() {
			playerName := strings.TrimSpace(nameInput.GetText())
			if len(playerName) > 0 {
				_, roomSelected := roomSelection.GetCurrentOption()
				if len(roomSelected) > 0 {
					// existing room
					if roomSelected != openRoomLabel {
						roomUrl = serverUrl + "/room-" + roomSelected + "/?arg=-name&arg=" + playerName
						flex.Clear()
						flex.AddItem(textView, 6, 1, false).AddItem(displayUrl, 10, 1, true).AddItem(githubLink, 1, 1, false)
					} else {
						// open new room
						newRoomName := strings.TrimSpace(newRoom.GetText())
						if len(newRoomName) > 0 {
							AddRoom(newRoomName)
							roomUrl = serverUrl + "/room-" + newRoomName + "/?arg=-name&arg=" + playerName
							flex.Clear()
							flex.AddItem(textView, 6, 1, false).
								AddItem(displayUrl, 10, 1, true).
								AddItem(githubLink, 1, 1, false)
						}
					}
					// display room url
					displayUrl.SetText(roomUrl)
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

	// Build screen
	flex.AddItem(textView, 6, 1, false).
		AddItem(form, 9, 1, true).
		AddItem(githubLink, 1, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	// refresh rooms available
	for {
		time.Sleep(time.Second * 5)

		rooms, err = RoomDeployed()
		if err != nil {
			log.Error().Msg("unable to get list of rooms deployed...")
		} else {
			log.Debug().Msgf("Found rooms: %v", rooms)
		}
		rooms = append(rooms, openRoomLabel)

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
	}
}