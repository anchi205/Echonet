package main

import (
	"flag"
	"fmt"
	"os" // way to interact with the OS
	"time"

	"github.com/anchi205/Echonet/src"
	"github.com/sirupsen/logrus" // structured logging.
)

// Called automatically before the main function
func init() {
	// Sets the log format to a colored, timestamped text format
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	// Log output to os.Stdout (standard output)
	logrus.SetOutput(os.Stdout)
}

func main() {
	// Define input flags
	username := flag.String("user", "", "Username to use in the chatroom.")
	chatroom := flag.String("room", "", "Chatroom to join.")
	loglevel := flag.String("log", "", "Level of logs to print.")
	discovery := flag.String("discover", "", "Method to use for discovery.")

	// Parse input flags
	flag.Parse()

	// Set the log level
	switch *loglevel {
	case "panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Display the introductory welcome msg
	fmt.Println("Echonet is starting.")
	fmt.Println("This may take a few seconds.")
	fmt.Println()

	// Create a new P2PHost
	p2phost := src.NewP2P()
	logrus.Infoln("P2P Setup Completed!")

	// Connect to peers with the chosen discovery method
	switch *discovery {
	case "announce":
		p2phost.AnnounceConnect()
		break
	case "advertise":
		p2phost.AdvertiseConnect()
		break
	default:
		p2phost.AdvertiseConnect()
	}

	logrus.Infoln("Connected to Peers")

	// Join the chat room
	chatapp, _ := src.JoinChatRoom(p2phost, *username, *chatroom)
	logrus.Infof("Joined the '%s' chatroom as '%s'", chatapp.RoomName, chatapp.UserName)

	// Wait for network setup to complete
	time.Sleep(time.Second * 5)

	// Create the Chat UI
	ui := src.NewUI()

	// Start the UI system
	ui.Run()
}
