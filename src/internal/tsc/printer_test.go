// +build windows
// +build amd64

package tsc

import (
	"log"
	"math"
	"testing"

	"git.xfyun.cn/hxzhao/am_sign/pkg/printer"
)

var _printer printer.Printer

func TestMain(m *testing.M) {
	var err error
	c := Config{
		Port:     "Gprinter GP-3120TU",
		Width:    "70",
		Height:   "50",
		Speed:    "5",
		Density:  "8",
		Sensor:   "0",
		Vertical: "1.5",
		Offset:   "0",
	}
	_printer, err = newTscPrinter(c)
	if err != nil {
		log.Fatalf("init printer get error: %s", err)
	}
	m.Run()
	if err := _printer.Print("1"); err != nil {
		log.Printf("commit print failed: %s", err)
	}
	if err := _printer.Close(); err != nil {
		log.Printf("close get error: %s", err)
	}
}
func TestTscPrinter_BarCode(t *testing.T) {
	if err := _printer.BarCode("300", "100", "128", "50", "1", "0", "2", "2", "123456789"); err != nil {
		t.Errorf("barcode failed: %s", err)
	}
}

func TestTscPrinter_QrCode(t *testing.T) {
	if err := _printer.QrCode("60", "50", "L", "6", "A", "0", "M2", "S3", "18119653669"); err != nil {
		t.Errorf("qrcode failed: %s", err)
	}
}

func TestTscPrinter_WindowsFont(t *testing.T) {
	content := "测试 windows font"
	if err := _printer.WindowsFont(calName(5, 90), 54, 90, 0, 2, 0, "Microsoft YaHei", content); err != nil {
		t.Errorf("windows font failed: %s", err)
	}
}
func TestTscPrinter_PrinterFont(t *testing.T) {
	content := "测试"

	if err := _printer.PrinterFont("300", "70", "3", "0", "1", "1", content); err != nil {
		t.Errorf("printer font failed %s", err)
	}
}

func calName(num float64, h float64) int {
	a := math.Round((600 - num*h*7/8.0 - h/8.0) / 2.0)
	if a < 0.0 {
		a = math.Round((600 - num*h*7/16.0 - h/8.0) / 2.0)
	}
	return int(a)
}
