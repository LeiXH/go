package printer

import (
	"fmt"
	"pkg/logger"

	"pkg/models"
)

const (
	nameYMul  = 45
	titleYMul = 140
)

func PrintUserLabel(userInfo *models.UserInfo) {
	doPrint(userInfo)
}

func doPrint(userInfo *models.UserInfo) {

	area := fmt.Sprintf("座位区域: %s", userInfo.Degree)
	printSignCard(userInfo.Name, userInfo.Department, userInfo.QrCode, area)

}

func printSignCard(name, title, code, area string) {

	pa := printArgs{Name: name, Title: title}
	if _printer == nil {
		logger.Error("你妈海, 竟然空指针")
		return
	}
	_ = _printer.WindowsFont(pa.getNameStartPoint(), nameYMul, pa.getNameHeight(), 0, 2, 0, "Microsoft YaHei", pa.getPrintableName())

	if len(title) > 0 {
		_ = _printer.WindowsFont(pa.getTitleStartPoint(), titleYMul, pa.getTitleHeight(), 0, 0, 0, "Microsoft YaHei", pa.getPrintableTitle())
	}

	_ = _printer.QrCode("207", "200", "L", "6", "A", "0", "M2", "S3", code)

	if len(area) > 0 {
		_ = _printer.WindowsFont(20, 350, 40, 0, 0, 0, "Microsoft YaHei", area)
	}

	_ = _printer.Print("1")
}
