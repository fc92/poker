package display

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"

	co "github.com/fc92/poker/internal/common"
)

func Display(localPlayer *co.Participant, room *co.Room, controlFromUI chan<- string, displayControl <-chan bool) {
	app := tview.NewApplication()
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	main := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetScrollable(false).
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	main.SetBorder(true).SetTitle("Team votes")
	voteList := tview.NewList().ShowSecondaryText(false)
	voteList.SetBorder(true).SetTitle("Player vote")
	commandList := tview.NewForm().SetButtonsAlign(tview.AlignCenter)

	gauge := tvxwidgets.NewPercentageModeGauge()
	gauge.SetMaxValue(1)
	gauge.SetTitle("% of votes received")
	gauge.SetRect(10, 4, 50, 3)
	gauge.SetBorder(true)

	barGraph := tvxwidgets.NewBarChart()
	barGraph.SetRect(4, 2, 50, 20)
	barGraph.SetBorder(true)
	barGraph.SetTitle("Distribution of votes")
	barGraph.AddBar(co.CommandVote1, 0, tcell.ColorGreen)
	barGraph.AddBar(co.CommandVote2, 0, tcell.ColorYellow)
	barGraph.AddBar(co.CommandVote3, 0, tcell.ColorBlue)
	barGraph.AddBar(co.CommandVote5, 0, tcell.ColorOrange)
	barGraph.AddBar(co.CommandVote8, 0, tcell.ColorIndianRed)
	barGraph.AddBar("13", 0, tcell.ColorRed)
	barGraph.AddBar("?", 0, tcell.ColorGray)

	grid := tview.NewGrid().
		SetRows(3, 0, 3, 3).
		SetColumns(15, 0, 40).
		AddItem(newPrimitive("Team poker\nPlayer: "+localPlayer.Name+"\n"), 0, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(voteList, 1, 0, 2, 1, 0, 100, true).
		AddItem(main, 1, 1, 2, 1, 0, 100, false).
		AddItem(barGraph, 1, 2, 1, 1, 0, 100, false).
		AddItem(gauge, 2, 2, 1, 1, 0, 100, false).
		AddItem(commandList, 3, 0, 1, 3, 0, 100, false)

	// wait for server update before refreshing
	go refreshLoop(displayControl, main, room, voteList, commandList, barGraph, gauge, localPlayer, controlFromUI, grid, app)

	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func refreshLoop(displayControl <-chan bool, main *tview.TextView, room *co.Room,
	voteList *tview.List, commandList *tview.Form, barGraph *tvxwidgets.BarChart,
	gauge *tvxwidgets.PercentageModeGauge, voter *co.Participant,
	controlFromUI chan<- string, grid *tview.Grid, app *tview.Application) {
	for {
		<-displayControl
		// update voter status
		main.SetText("")
		fmt.Fprintf(main, "%s", strings.Join(room.DisplayVotersStatus(), ""))

		// update command available
		updateCommands(voteList, commandList, voter, room, controlFromUI, app)
		// update bar chart
		updateBarChart(barGraph, *room, grid)
		// update gauge
		updateGauge(gauge, *room)
		app.Draw().SetFocus(voteList)
	}
}

func updateCommands(voteList *tview.List, commandList *tview.Form, voter *co.Participant, room *co.Room, controlFromUI chan<- string, app *tview.Application) {
	voteList.Clear()
	commandList.ClearButtons()
	commands := []string{}
	for k := range voter.AvailableCommands {
		commands = append(commands, k)
	}
	sort.Strings(commands)
	for _, shortcut := range commands {
		key := []rune(shortcut)[0]
		shortcut := shortcut // need to differentiate the shortcut variables to differentiate anonymous func
		// handle vote commands
		if _, ok := room.VoteCommands()[shortcut]; ok {
			voteList.AddItem(voter.AvailableCommands[shortcut], "", key, func() {
				controlFromUI <- shortcut
			})
			// highlight current vote
			if voter.Vote == voter.AvailableCommands[shortcut] {
				voteList.SetCurrentItem(voteList.FindItems(voter.Vote, "", false, false)[0])
				voteList.SetSelectedTextColor(tcell.ColorRed)
				voteList.SetChangedFunc(func(ind int, mText string, sText string, scut rune) {
					if voter.Vote == voter.AvailableCommands[string(scut)] {
						voteList.SetSelectedTextColor(tcell.ColorRed)
					} else {
						voteList.SetSelectedTextColor(tcell.ColorBlack)
					}
				})
			}
		} else {
			// handle other commands
			commandList.AddButton(voter.AvailableCommands[shortcut], func() {
				controlFromUI <- shortcut
			})

		}
	}
	// reset default selection if no local vote is available
	if voter.Vote == co.VoteNotReceived {
		voteList.SetCurrentItem(0)
		voteList.SetSelectedTextColor(tcell.ColorBlack)
	}
	commandList.AddButton("Quit", func() {
		app.Stop()
		controlFromUI <- co.CommandQuit
	})
}

func updateBarChart(barGraph *tvxwidgets.BarChart, room co.Room, grid *tview.Grid) {
	if room.RoomStatus == co.VoteClosed {
		grid.AddItem(barGraph, 1, 2, 1, 1, 0, 100, false)
		voteSum := map[string]int{}
		for _, voter := range room.Voters {
			voteSum[voter.Vote]++
		}

		barGraph.SetBarValue(co.CommandVote1, voteSum[room.VoteCommands()[co.CommandVote1]])
		barGraph.SetBarValue(co.CommandVote2, voteSum[room.VoteCommands()[co.CommandVote2]])
		barGraph.SetBarValue(co.CommandVote3, voteSum[room.VoteCommands()[co.CommandVote3]])
		barGraph.SetBarValue(co.CommandVote5, voteSum[room.VoteCommands()[co.CommandVote5]])
		barGraph.SetBarValue(co.CommandVote8, voteSum[room.VoteCommands()[co.CommandVote8]])
		barGraph.SetBarValue("13", voteSum[room.VoteCommands()[co.CommandVote13]])
		barGraph.SetBarValue("?", voteSum[room.VoteCommands()[co.CommandNotVoting]])
		barGraph.SetMaxValue(len(room.Voters))
	} else {
		barGraph.SetBarValue(co.CommandVote1, 0)
		barGraph.SetBarValue(co.CommandVote2, 0)
		barGraph.SetBarValue(co.CommandVote3, 0)
		barGraph.SetBarValue(co.CommandVote5, 0)
		barGraph.SetBarValue(co.CommandVote8, 0)
		barGraph.SetBarValue("13", 0)
		barGraph.SetBarValue("?", 0)
		barGraph.SetMaxValue(len(room.Voters))
	}
}

func updateGauge(gauge *tvxwidgets.PercentageModeGauge, room co.Room) {
	gauge.SetMaxValue(len(room.Voters))
	gauge.SetValue(room.NbVotesReceived())
}
