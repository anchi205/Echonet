package src

// Step 1: Import Packages
import (
	"fmt"
	"io"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Step 2: Components and channels needed for the UI and communication.
type UI struct {
	TerminalApp *tview.Application // Instance of the tview.Application for managing the terminal UI.
	PeerBox     *tview.Box         // displaying information about connected peers.
	ChatBox     io.Writer          // handling messages to be displayed in the chat room.

	LogChan   chan string   // Channel for receiving log messages.
	InputChan chan string   // Channel for receiving user input messages.
	SyncChan  chan struct{} // Channel for signaling synchronization events.
	TermChan  chan struct{} // Channel for signaling termination events.
}

func NewUI() *UI {
	/*
			Initializes and configures the UI components.
		Creates text views for the title, chat room, usage information, peers, and user input.
		Sets up an input field for users to type messages.
		Defines callbacks for handling input, including commands like /quit, /sync, /room, and /user.
		Creates a flexible layout using tview.NewFlex to arrange UI components in rows and columns.
		Sets the root of the UI with the configured layout.
	*/

	// Step 3: Initialize the tview app
	app := tview.NewApplication()

	syncchan := make(chan struct{})
	inputchan := make(chan string)
	logchan := make(chan string)

	// Step 4: UI components
	commands := tview.NewTextView().SetDynamicColors(true).SetText(`
					[red]/quit[green] - Exit the chat |
					[red]/room <roomname>[green] - Change chat room |
					[red]/user <username>[green] - Change user name |
					[red]/sync[green] - refresh`)

	commands.
		SetBorder(true).
		SetBorderColor(tcell.ColorRebeccaPurple).
		SetTitle("Commands").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite).
		SetBorderPadding(0, 0, 2, 0)

	titlebox := tview.NewTextView().
		SetText("Echonet. A Golang based P2P Chat Application. ").
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignCenter)

	titlebox.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen)

	chatbox := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	chatbox.
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle(fmt.Sprintf("ChatRoom-%s", defaultroom)).
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite)

	peerbox := tview.NewTextView().
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle("Peers").
		SetTitleAlign(tview.AlignRight).
		SetTitleColor(tcell.ColorWhite)

	input := tview.NewInputField().
		SetLabel(defaultuser + " > ").
		SetLabelColor(tcell.ColorGreen).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitle("Input").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite).
		SetBorderPadding(0, 0, 1, 0)

		// Step 5: Handle Input
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter { // Check if trigger was caused by a Return(Enter) press.
			return
		}

		line := input.GetText()

		if len(line) == 0 {
			return
		}

		// Check for command inputs
		if strings.HasPrefix(line, "/") {
			if strings.HasPrefix(line, "/quit") {
				app.Stop()
				return
			} else if strings.HasPrefix(line, "/sync") {
				syncchan <- struct{}{}

			} else if strings.HasPrefix(line, "/room") {
				// room change cmd
			} else if strings.HasPrefix(line, "/user") {
				// user change cmd
			} else {
				logchan <- "Error. invalid command!"
			}
		}

		// Send the message to the input channel
		inputchan <- line

		// Reset the input field
		input.SetText("")
	})

	// Step 6: Create Layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titlebox, 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(chatbox, 0, 1, false).
			AddItem(peerbox, 20, 1, false),
			0, 1, false).
		AddItem(input, 3, 1, false).
		AddItem(commands, 3, 1, false)

	app.SetRoot(flex, true)

	return &UI{
		TerminalApp: app,
		PeerBox:     peerbox,
		ChatBox:     chatbox,
		TermChan:    make(chan struct{}, 1),
	}
}

// Starts the UI application and runs it until it's stopped.
func (ui *UI) Run() error {
	defer ui.Close()
	return ui.TerminalApp.Run()
}

// Signals the termination of the application by sending a message to TermChan
func (ui *UI) Close() {
	ui.TermChan <- struct{}{}
}
