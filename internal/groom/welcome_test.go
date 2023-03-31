// package groom is used to add or remove poker room
// based on a Welcome screen in terminal mode
package groom

import (
	"reflect"
	"testing"

	"github.com/rivo/tview"
)

func Test_projectLink(t *testing.T) {
	type args struct {
		app *tview.Application
	}
	tests := []struct {
		name string
		args args
		want *tview.TextView
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := projectLink(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("projectLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setRoomSelectionOptions(t *testing.T) {
	type args struct {
		roomSelection *tview.DropDown
		form          *tview.Form
		newRoom       *tview.InputField
		app           *tview.Application
		openRoomLabel string
		newRoomLabel  string
		rooms         *[]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// add 3 rooms
		{
			name: "3 rooms",
			args: args{
				roomSelection: tview.NewDropDown(),
				form:          tview.NewForm(),
				newRoom:       tview.NewInputField(),
				app:           tview.NewApplication(),
				openRoomLabel: "Open a room",
				newRoomLabel:  "Create a new room",
				rooms: &[]interface{}{
					map[string]interface{}{
						"name": "room1",
					},
					map[string]interface{}{
						"name": "room2",
					},
					map[string]interface{}{
						"name": "Open a room",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRoomSelectionOptions(tt.args.roomSelection, tt.args.form, tt.args.newRoom, tt.args.app, tt.args.openRoomLabel, tt.args.newRoomLabel, tt.args.rooms)
		})
	}
}

func Test_getRoomsName(t *testing.T) {
	type args struct {
		rooms []interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// add 3 rooms based on values.yaml syntax
		{
			name: "3 rooms",
			args: args{
				rooms: []interface{}{
					map[string]interface{}{
						"name": "room1",
					},
					map[string]interface{}{
						"name": "room2",
					},
					map[string]interface{}{
						"name": "room3",
					},
				},
			},
			want: []string{"room1", "room2", "room3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRoomsName(tt.args.rooms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRoomsName() = %v, want %v", got, tt.want)
			}
		})
	}
}
