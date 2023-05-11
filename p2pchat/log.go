package main

import (
	"bytes"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// go-ethereum doesn't provide custom Logger factory (log.logger is private)
// Consider replace with other logging library
// Here we expose geth's default logger and set level to "WARN"
func init() {
	logger := log.Root()
	geth_formatter := log.TerminalFormat(false)
	formatter := log.FormatFunc(func(r *log.Record) []byte {
		buf := &bytes.Buffer{}
		// buf.WriteString(fmt.Sprintln(
		// 	r.Time.Format(time.RFC3339),
		// 	fmt.Sprintf("[%v]", r.Lvl.AlignedString()),
		// 	r.Msg,
		// ))
		buf.WriteString("Geth ")
		buf.Write(geth_formatter.Format(r))
		return buf.Bytes()
	})
	handler := &log.GlogHandler{}
	// set reasonable defaults
	handler.SetHandler(log.StreamHandler(os.Stdout, formatter))
	handler.Verbosity(log.LvlWarn)
	logger.SetHandler(handler)
}

// NOTE: This logger is currently unavailable because it is highly-customization
// for go-ethereum.
// func GetLogger() log.Logger { return logger }

func SetLogLevel(level log.Lvl) {
	handler := log.Root().GetHandler().(*log.GlogHandler)
	handler.Verbosity(level)
}

// NOTE: geth log ctx is key,value mapping for logging
// log.Error("Could not search for pattern", "pattern", pattern, "contract", contracts[types[i]], "err", err)
