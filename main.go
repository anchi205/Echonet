package main

import (
	"os" // way to interact with the OS, such as reading envt variables.

	"github.com/anchi205/Echonet/src"
	"github.com/sirupsen/logrus" // structured logging.
)

func init() { // called automatically before the main function.
	// It sets the log format to a colored, timestamped text format
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// Configures log output to os.Stdout (standard output).
	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)

	// Sets the log level to InfoLevel - log informational warning msgs and above.
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	ui := src.NewUI()
	ui.Run()
}
