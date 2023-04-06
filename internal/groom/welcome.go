// package groom is used to add or remove poker room
// based on a Welcome screen in terminal mode
package groom

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/common-nighthawk/go-figure"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common/logger"
)

const (
	newRoomLabel  = "New room name (4-10 letters only):"
	openRoomLabel = "open new room"
	urlLabel      = "Use this url to join the poker room : "
	errorLabel    = "Room is not available, please try again later"
	tipsLabel     = "(pop-up needs to be allowed, multiple click work best)"
	inputSize     = 20
)

func init() {
	logger.InitLogger()
}

// Terminal welcome page to choose player name and poker room
// serverUrl designate the url used by the browser to reach the poker room
func DisplayWelcome(serverURL string) {
	rooms, err := RoomDeployed()
	if err != nil {
		log.Err(err).Msg("unable to get list of rooms deployed...")
	} else {
		log.Debug().Msgf("Found initial rooms: %v", rooms)
	}
	rooms = append(rooms, map[string]interface{}{"name": openRoomLabel, "index": -1})
	app := tview.NewApplication()
	roomURL := ""
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Title
	pokerFigure := figure.NewFigure("Team Poker", "", true)
	title := pokerFigure.String()

	// github link
	githubLink := projectLink(app)

	// Form fields
	titleView := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).
		SetTextAlign(0).
		SetMaxLines(10).
		SetText(title).
		SetChangedFunc(func() {
			app.Draw()
		})
	nameInput := tview.NewInputField().SetLabel("Enter your name:").SetFieldWidth(inputSize)
	newRoom := tview.NewInputField().SetLabel(newRoomLabel).SetFieldWidth(inputSize)
	newRoom.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {

		// each char must be a letter
		for _, char := range textToCheck {
			if !unicode.IsLetter(char) {
				return false
			}
		}

		return true
	})
	displayURL := tview.NewTextView().SetLabel(urlLabel)
	tipsURL := tview.NewTextView().SetLabel(tipsLabel)
	roomSelection := tview.NewDropDown().
		SetFieldWidth(inputSize).
		SetLabel("Select poker room:")

	// Form
	form := newForm(nameInput, roomSelection, roomURL, serverURL, flex, titleView, displayURL, tipsURL, githubLink, newRoom, app)

	// Show/Hide new room name in form
	setRoomSelectionOptions(roomSelection, form, newRoom, app, openRoomLabel, newRoomLabel, &rooms)

	// Build screen
	flex.AddItem(titleView, 6, 1, false).
		AddItem(form, 9, 1, true).
		AddItem(githubLink, 1, 1, false)

	go app.SetRoot(flex, true).EnableMouse(true).Run()

	// refresh rooms available
	refreshRoomList(rooms, err, roomSelection, form, newRoom, app)
}

func refreshRoomList(rooms []interface{}, err error, roomSelection *tview.DropDown, form *tview.Form, newRoom *tview.InputField, app *tview.Application) {
	for {
		log.Debug().Msg("starting room list refresh")
		time.Sleep(time.Second * 5)

		var err2 error
		rooms, err2 = RoomDeployed()
		if err2 != nil {
			log.Err(err2).Msg("unable to get list of rooms deployed...")
			return
		} else {
			log.Debug().Msgf("Found rooms: %v", rooms)
		}
		rooms = append(rooms, map[string]interface{}{"name": openRoomLabel, "index": -1})

		setRoomSelectionOptions(roomSelection, form, newRoom, app, openRoomLabel, newRoomLabel, &rooms)
		app.Draw()
	}
}

func projectLink(app *tview.Application) *tview.TextView {
	githubLink := tview.NewTextView().SetTextColor(tview.Styles.PrimaryTextColor).
		SetTextAlign(tview.AlignRight).
		SetMaxLines(1).
		SetText("https://github.com/fc92/poker").
		SetChangedFunc(func() {
			app.Draw()
		})
	return githubLink
}

func newForm(nameInput *tview.InputField, roomSelection *tview.DropDown, roomURL string, serverURL string, flex *tview.Flex, textView *tview.TextView, displayURL *tview.TextView, tipsURL *tview.TextView, githubLink *tview.TextView, newRoom *tview.InputField, app *tview.Application) *tview.Form {
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
						roomURL = serverURL + "/room-" + roomSelected + "/?arg=-name&arg=" + playerName
						displayURL.SetText(roomURL)
						displayResultURL(flex, textView, displayURL, tipsURL, githubLink)
					} else {
						// open new room
						newRoomName := strings.TrimSpace(newRoom.GetText())
						if len(newRoomName) > 3 && len(newRoomName) < 11 {
							AddRoom(newRoomName)
							roomURL = serverURL + "/room-" + newRoomName + "/?arg=-name&arg=" + playerName
							displayURL.SetText(roomURL)
							displayResultURL(flex, textView, displayURL, tipsURL, githubLink)
						}
					}
				}
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
			fmt.Print("Exiting poker groom. You can close this web page or refresh it to join again.")
			os.Exit(0)
		})
	return form
}

func setRoomSelectionOptions(roomSelection *tview.DropDown, form *tview.Form, newRoom *tview.InputField, app *tview.Application, openRoomLabel string, newRoomLabel string, rooms *[]interface{}) {
	roomSelection.SetOptions(getRoomsName(*rooms), func(option string, index int) {
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

func displayResultURL(flex *tview.Flex, titleView *tview.TextView, displayURL *tview.TextView, tipsURL *tview.TextView, githubLink *tview.TextView) {
	// textView to display waitLabel
	waitView := tview.NewTextView().SetText("Checking room availability...")

	// display waitView
	flex.Clear()
	flex.AddItem(titleView, 6, 1, false).
		AddItem(waitView, 10, 1, true).
		AddItem(githubLink, 1, 1, false)

	// during 60 s check if url is reachable every 2 seconds
	go func() {
		for i := 0; i < 30; i++ {
			if checkURL(displayURL.GetText(false)) {
				flex.Clear()
				flex.AddItem(titleView, 6, 1, false).
					AddItem(displayURL, 10, 1, true).
					AddItem(tipsURL, 10, 1, true).
					AddItem(githubLink, 1, 1, false)
				break
			}
			time.Sleep(time.Second * 2)
			waitView.SetText("Checking room availability... " + strconv.Itoa((i+1)*2) + " seconds")
			if i == 29 {
				// textView to display errorLabel
				errorView := tview.NewTextView().SetText(errorLabel)
				// display errorView
				flex.Clear()
				flex.AddItem(titleView, 6, 1, false).
					AddItem(errorView, 10, 1, true).
					AddItem(githubLink, 1, 1, false)
			}
		}
	}()

}

func getRoomsName(rooms []interface{}) []string {
	var options []string
	for _, room := range rooms {
		if roomMap, ok := room.(map[string]interface{}); ok {
			if name, ok := roomMap["name"].(string); ok {
				options = append(options, name)
			}
		}
	}
	return options
}

// function to check if url is reachable with http code 200
func checkURL(url string) bool {
	// clean url by removing arguments
	url = strings.Split(url, "?")[0]
	response, err := http.Get(url) // do not check certificate
	if err != nil {
		log.Warn().Msgf("Error while checking url %s: %s", url, err)
		return false
	}
	//if http code is 503 or 404 return false
	if response.StatusCode == 503 || response.StatusCode == 404 {
		log.Info().Msgf("Url %s is not reachable", url)
		return false
	}

	return true
}
