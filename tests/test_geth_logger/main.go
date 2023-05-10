package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// go-ethereum doesn't provide custom Logger factory (log.logger is private)
// Consider replace with other logging library
var logger = log.Root()

func init() {
	formatter := log.LogfmtFormat()
	handler := &log.GlogHandler{}
	// set reasonable defaults
	handler.SetHandler(log.StreamHandler(os.Stdout, formatter))
	handler.Verbosity(log.LvlWarn)
	logger.SetHandler(handler)
}

func setLogLevel(level log.Lvl) {
	handler := logger.GetHandler().(*log.GlogHandler)
	handler.Verbosity(level)
}

func doLog(log log.Logger) {
	log.Trace("Trace")
	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
}

func main() {
	doLog(logger)
	fmt.Println("changed log level")
	setLogLevel(log.LvlDebug)
	doLog(logger)
}
