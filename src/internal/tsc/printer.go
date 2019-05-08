// +build windows
// +build amd64

package tsc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pkg/logger"
	"strings"
	"syscall"
)

type tscPrinter struct {
	lib syscall.Handle

	functions map[string]uintptr

	opened bool

	config Config
}

func newTscPrinter(c Config) (*tscPrinter, error) {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}

	path := filepath.Join(pwd, "./bin/TSCLIB.dll")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Infof("err %s", pwd)
		logger.Warnf("pwd get error, dll load maybe failed")
		path = "TSCLIB.dll"
	}
	tp := &tscPrinter{config: c}
	tp.functions = make(map[string]uintptr)

	if err := tp.loadLibrary(path); err != nil {
		return nil, err
	}
	if err := tp.openPort(tp.config.Port); err != nil {
		return nil, err
	}
	tp.opened = true
	if err := tp.setup(tp.config.Width, tp.config.Height, tp.config.Speed, tp.config.Density, tp.config.Sensor, tp.config.Vertical, tp.config.Offset); err != nil {
		return nil, err
	}
	if err := tp.clearBuffer(); err != nil {
		return nil, err
	}
	return tp, nil
}

// Initialization of Printer
// if the dynamic library is not proper, will return an error
// but if function in the library not proper, will panic
func (p *tscPrinter) loadLibrary(path string) error {
	var err error
	p.lib, err = syscall.LoadLibrary(path)
	if err != nil {
		return fmt.Errorf("load library %s got error: %s", path, err)
	}
	return p.preFlight()
}

// free the printer
// close port using and free the dynamic library
func (p *tscPrinter) Close() error {
	if err := p.clearBuffer(); err != nil {
		return err
	}
	if err := p.closePort(); err != nil {
		return err
	}
	p.opened = false
	return syscall.FreeLibrary(p.lib)
}

func (p *tscPrinter) preFlight() (err error) {
	for _, procName := range []string{"about", "openport", "closeport", "setup", "clearbuffer", "barcode",
		"printerfont", "sendcommand", "downloadpcx", "printlabel", "nobackfeed", "windowsfont", "formfeed"} {
		p.functions[procName], err = findProcAddr(p.lib, procName)
		if err != nil {
			return fmt.Errorf("tsc printer preflight failed, can't load function %s, because: %s", procName, err)
		}
	}
	return
}

func (p *tscPrinter) preCheck() (err error) {
	if !p.opened {
		return errors.New("printer have been closed")
	}
	return nil
}

func (p *tscPrinter) getFunction(name string) uintptr {
	_func := p.functions[name]
	return _func
}

// Start the windows printer spool
// for local printer, port is the printer diver name, like "TTP-244 Plus"
// for network printer port is the UNC path and printer name, like "\\server\TTP243"
// for centronics interface directly, please specify LPT1 to LPT4, like "LPT1"
// for USB interface directly, please specify USB, like "USB"
func (p *tscPrinter) openPort(port string) error {
	_func := p.getFunction("openport")

	var nargs uintptr = 1
	_, _, callErr := syscall.Syscall(_func, nargs,
		stringToUintptr(port),
		0,
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function openport got error, port is %s, err is %s", port, callErr)
	}
	return nil
}

// close windows printer spool
func (p *tscPrinter) closePort() error {
	if !p.opened {
		return nil
	}
	_func := p.getFunction("closeport")

	var nargs uintptr = 0
	_, _, callErr := syscall.Syscall(_func, nargs, 0, 0, 0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function closeport got error %s", callErr)
	}
	return nil
}

// Setup label width, label height, print speed, print density, snesor type, gap/black mark vertical distance, gap/black markshift distance
// width and heigth unit: mm
// selectable print speeds vary on different printer models, the valid values are 1.0, 1.5, 2.0, 3.0, 4.0, 6.0, 8.0, 10.0, 12.0,  unit "/sec
// density, the greater the number, the darker the printing, 0~15 are valid
// sensor, 0 signifies the vertical gap sensor is to be used; 1 signifies that black mark sensor is to be used
// vertical: the height of gap/the vertical distance of black mark, unit mm
// offset: gap/black markshift distance, unit: mm, in the case of the average labe, set this parameter to be 0
func (p *tscPrinter) setup(width, height, speed, density, sensor, vertical, offset string) error {
	if !p.opened {
		return errors.New("printer have been closed")
	}

	_func := p.getFunction("setup")

	var nargs uintptr = 8
	_, _, callErr := syscall.Syscall9(_func, nargs,
		stringToUintptr(width),
		stringToUintptr(height),
		stringToUintptr(speed),
		stringToUintptr(density),
		stringToUintptr(sensor),
		stringToUintptr(vertical),
		stringToUintptr(vertical),
		stringToUintptr(offset),
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function setup got error %s", callErr)
	}
	return nil
}

// clear, like reset
func (p *tscPrinter) clearBuffer() error {
	if !p.opened {
		return nil
	}
	_func := p.getFunction("clearbuffer")

	var nargs uintptr = 0
	_, _, callErr := syscall.Syscall(_func, nargs, 0, 0, 0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function clearbuffer got error %s", callErr)
	}
	return nil
}

// use built-in bar code formats to print
// the starting point of the bar code along the X/Y direction, given in points, (of 200DPI, 1 point = 1/8 mm; of 300DPI 1 point = 1/12 mm)
// _type: barcode style, see TSC_DLL_instruction.pdf for detail, https://en.wikipedia.org/wiki/Barcode#Linear_barcodes
// height: barcode height, in points
// readable: whether to print human recognizable interpretation(text), 0 not, 1 will
// rotation degrees, 0, 90, 180, or 270 in counter clockwise direction
//  narrow, wide : narrow bar ratio, see TSC_Program_Manual.pdf for detail
// code: content of barcode, ABC-abc-1234
func (p *tscPrinter) BarCode(x, y, _type, height, readable, rotation, narrow, wide, code string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("barcode")

	var nargs uintptr = 9
	_, _, callErr := syscall.Syscall9(_func, nargs,
		stringToUintptr(x),
		stringToUintptr(y),
		stringToUintptr(_type),
		stringToUintptr(height),
		stringToUintptr(readable),
		stringToUintptr(rotation),
		stringToUintptr(narrow),
		stringToUintptr(wide),
		stringToUintptr(code))
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function barcode got error %s", callErr)
	}
	return nil
}

// use printer built-in fonts to print
// the starting point of the bar code along the X/Y direction, given in points, (of 200DPI, 1 point = 1/8 mm; of 300DPI 1 point = 1/12 mm)
// fonttype: built-in font type name, see TSC_DLL_instruction.pdf for detail
// rotation degrees, 0, 90, 180, or 270 in counter clockwise direction
// sets up the magnification rate of the content along the X/Y direction, range 1~8
// the text to print, in utf8-encoding, this function will transform it to GBK
func (p *tscPrinter) PrinterFont(x, y, fontType, rotation, xMul, yMul, content string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("printerfont")

	// TODO
	gbkContent, err := utf8ToGBK([]byte(content))
	if err != nil {
		return fmt.Errorf("tsc printer: tansform content %s to gbk encode failed", content)
	}

	var nargs uintptr = 7
	_, _, callErr := syscall.Syscall9(_func, nargs,
		stringToUintptr(x),
		stringToUintptr(y),
		stringToUintptr(fontType),
		stringToUintptr(rotation),
		stringToUintptr(xMul),
		stringToUintptr(yMul),
		bytesToUintptr(gbkContent),
		0,
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function printerfont got error %s", callErr)
	}
	return nil
}

func (p *tscPrinter) QrCode(x, y, eccLevel, cellWidth, mode, rotation, model, mask string, content string) error {
	base := "QRCODE"
	cmd := fmt.Sprintf("%s %s,\"%s\"", base, strings.Join([]string{x, y, eccLevel, cellWidth, mode, rotation, model, mask}, ","), content)
	return p.sendCommand(cmd)
}

func (p *tscPrinter) Print(_copy string) error {

	err := p.printLabel("1", _copy)
	if err != nil {
		return err
	}
	if err := p.clearBuffer(); err != nil {
		return err
	}
	if err := p.closePort(); err != nil {
		return err
	}
	if err := p.openPort(p.config.Port); err != nil {
		return err
	}
	return nil
}

// use windows font to print text
// the starting point of the bar code along the X/Y direction, given in points, (of 200DPI, 1 point = 1/8 mm; of 300DPI 1 point = 1/12 mm)
// fontheight: the font height, given in points
// rotation degrees, 0, 90, 180, or 270 in counter clockwise direction
// fontstyle 0:Normal, 1: Italic, 2: Bold, 3: Bold and Intalic
// fontunderline: 0: without underline, 1: with underline
// szFaceName: specify the true type font name, e.g. Arial, Times new Roman, case sensitive
// the text to print, in utf8-encoding, this function will transform it to GBK
func (p *tscPrinter) WindowsFont(x, y, fontheight, rotation, fontstyle, fontunderline int, szFaceName, content string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("windowsfont")

	gbkContent, err := utf8ToGBK([]byte(content))
	if err != nil {
		return fmt.Errorf("tsc printer: tansform content %s to gbk encode failed", content)
	}

	var nargs uintptr = 8
	_, _, callErr := syscall.Syscall9(_func, nargs,
		uintptr(x),
		uintptr(y),
		uintptr(fontheight),
		uintptr(rotation),
		uintptr(fontstyle),
		uintptr(fontunderline),
		stringToUintptr(szFaceName),
		bytesToUintptr(gbkContent),
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function windowsfont got error %s", callErr)
	}
	return nil
}

// skip to next page(of label), this function is to be used after setup
func (p *tscPrinter) FormFeed() error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("formfeed")

	var nargs uintptr = 0
	_, _, callErr := syscall.Syscall(_func, nargs, 0, 0, 0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function formfeed got error %s", callErr)
	}
	return nil
}

// Display the DLL version on the screen
func (p *tscPrinter) about() error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("about")

	var nargs uintptr = 0
	_, _, callErr := syscall.Syscall(_func, nargs, 0, 0, 0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function about got error %s", callErr)
	}
	return nil
}

// download mono PCX graphi files to the printer
// filename, retrieval/absolute path
// imageName name of files that to be downloaded in the printer memory
// !!! Please use capital letters
func (p *tscPrinter) downloadPcx(filename, imageName string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("downloadpcx")

	var nargs uintptr = 2
	_, _, callErr := syscall.Syscall(_func, nargs,
		stringToUintptr(filename),
		stringToUintptr(imageName),
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function downloadpcx got error %s, filename is %s, imageName is %s", callErr, filename, imageName)
	}
	return nil
}

// disable the backfeed function
func (p *tscPrinter) noBackFeed() error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("nobackfeed")

	var nargs uintptr = 0
	_, _, callErr := syscall.Syscall(_func, nargs, 0, 0, 0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function nobackfeed got error %s", callErr)
	}
	return nil
}

// real to print label
// set up the label sets
// set up the number of print copies
func (p *tscPrinter) printLabel(set, _copy string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("printlabel")

	var nargs uintptr = 2
	_, _, callErr := syscall.Syscall(_func, nargs,
		stringToUintptr(set),
		stringToUintptr(_copy),
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function printlabel got error %s, set is %s, copy is %s", callErr, set, _copy)
	}
	return nil
}

// send built-in commands to the bar code printer
// see SC_Program_Manual.pdf for detail
func (p *tscPrinter) sendCommand(printerCommand string) error {
	if err := p.preCheck(); err != nil {
		return err
	}
	_func := p.getFunction("sendcommand")

	var nargs uintptr = 1
	_, _, callErr := syscall.Syscall(_func, nargs,
		stringToUintptr(printerCommand),
		0,
		0)
	if callErr != 0 {
		return fmt.Errorf("tsc printer: call function sendcommand got error %s, command is %s", callErr, printerCommand)
	}
	return nil
}
