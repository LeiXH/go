// +build windows
// +build amd64

package tsc

import (
	"encoding/json"
	"errors"
	"pkg/printer"
)

type PrinterFactory struct{}

// return a new tsc-Printer
func (pf PrinterFactory) NewPrinter(args string) (printer.Printer, error) {
	conf := Config{}
	if err := json.Unmarshal([]byte(args), &conf); err != nil {
		return nil, errors.New("parse printer config failed")
	}
	return newTscPrinter(conf)
}

func init() {
	printer.Register("tsc", PrinterFactory{})
}
