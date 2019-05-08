package printer

import (
	"encoding/json"
	"pkg/printer"
	"sync"

	_ "internal/tsc"

)

var _printer printer.Printer
var _printOnce sync.Once

func InitPrinter(config interface{}) (err error) {
	var cs []byte
	cs, err = json.Marshal(config)

	if err != nil {
		return err
	}
	_printOnce.Do(func() {
		_printer, err = printer.New("tsc", string(cs))
	})
	return
}
func ClosePrinter() {
	if _printer == nil {
		return
	}
	_ = _printer.Close()
}
