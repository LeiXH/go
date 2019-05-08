// +build windows
// +build amd64

package printer

import (
	"log"
	"testing"

	"pkg/models"
)
import _ "internal/tsc"

var printerCondig = struct {
	Port string

	Width    string
	Height   string
	Speed    string
	Density  string
	Sensor   string
	Vertical string
	Offset   string
}{
	Port:     "Gprinter GP-3120TU",
	Width:    "70",
	Height:   "50",
	Speed:    "5",
	Density:  "8",
	Sensor:   "0",
	Vertical: "1.5",
	Offset:   "0",
}

func TestMain(m *testing.M) {
	if err := InitPrinter(printerCondig); err != nil {
		log.Fatalf("init printer failed : %s", err)
	}
	defer ClosePrinter()
	m.Run()
}

func TestPrintUserLabel(t *testing.T) {
	PrintUserLabel(&models.UserInfo{Name: "何莎莎", WorkNumber: "00001", DomainName: "qfliu", Department: "消费者事业群", SeatArea: "看台 南C区01排001号", QrCode: "iflytek.com"})
}
