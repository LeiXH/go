package printer

import (
	"math"

	"pkg/utils"
)

const maxNameLength = 6
const namePrePrintHeight = 90
const titlePrePrintHeight = 40
const leastTitleStart = 20
const asciiStart = 15

type printArgs struct {
	Name  string
	Title string

	namePrintHeight  int
	titlePrintHeight int
}

func (pa *printArgs) getNameStartPoint() int {

	if utils.IfHaveAscii(pa.Name) {
		return asciiStart
	}

	ll := utils.CountUnicodeString(pa.Name)
	if ll > maxNameLength {
		ll = maxNameLength
	}
	pa.Name = utils.CutUnicodeString(pa.Name, ll)
	return int(math.Round((600 - float64(ll)*namePrePrintHeight*7/8.0 - namePrePrintHeight/8.0) / 2.0))
}
func (pa printArgs) getNameHeight() int {
	return namePrePrintHeight
}

func (pa printArgs) getPrintableName() string {
	return pa.Name
}

func (pa *printArgs) getTitleStartPoint() int {

	if utils.IfHaveAscii(pa.Title) {
		return asciiStart
	}

	ll := float64(utils.CountUnicodeString(pa.Title))
	_s := math.Round((600 - ll*titlePrePrintHeight*7/8.0 - titlePrePrintHeight/8.0) / 2.0)
	if _s <= 0 {
		pa.titlePrintHeight = int(math.Round(600-leastTitleStart) * 8 / (7*ll + 1))
		return leastTitleStart
	}
	return int(_s)
}

func (pa printArgs) getTitleHeight() int {
	if pa.titlePrintHeight <= 0 {
		return titlePrePrintHeight
	}
	return pa.titlePrintHeight
}

func (pa printArgs) getPrintableTitle() string {
	return pa.Title
}
