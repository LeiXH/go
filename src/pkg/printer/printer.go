package printer

import (
	"fmt"
	"sync"
)

var (
	printersMu sync.RWMutex
	printers   = make(map[string]Factory)
)

type Printer interface {
	BarCode(x, y, _type, height, readable, rotation, narrow, wide, code string) error
	QrCode(x, y, eccLevel, cellWidth, mode, rotation, model, mask string, content string) error
	PrinterFont(x, y, fontType, rotation, xMul, yMul, content string) error
	WindowsFont(x, y, fontHeight, rotation, fontStyle, fontUnderline int, szFaceName, content string) error
	FormFeed() error
	Print(_copy string) error
	Close() error
}

type Factory interface {
	NewPrinter(args string) (Printer, error)
}

func Register(name string, driver Factory) {
	printersMu.Lock()
	defer printersMu.Unlock()
	if driver == nil {
		panic("printer: Register Factory is nil")
	}
	if _, dup := printers[name]; dup {
		panic("printer: Register called twice for factory " + name)
	}
	printers[name] = driver
}

func New(name string, args string) (Printer, error) {
	fac, ok := printers[name]
	if !ok {
		return nil, fmt.Errorf("printer: no this printer")
	}
	return fac.NewPrinter(args)
}
