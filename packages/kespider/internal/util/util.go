package util

import (
	"regexp"
	"strings"
)

func ParseHouseDetail(houseInfo string) (string, string, string, string, string) {
	emptyReg, _ := regexp.Compile(`\s`)
	houseInfo = emptyReg.ReplaceAllString(houseInfo, "")
	houseTypeReg, _ := regexp.Compile(`\d+室\d+厅`)
	houseType := houseTypeReg.FindString(houseInfo)
	houseAreaReg, _ := regexp.Compile(`[1-9]\d*(\.\d*)?平米`)
	houseArea := houseAreaReg.FindString(houseInfo)
	houseOrientationReg, _ := regexp.Compile(`[东北西南]`)
	houseOrientation := strings.Join(houseOrientationReg.FindAllString(houseInfo, -1), "")
	houseYearReg, _ := regexp.Compile(`[1-9]\d*年`)
	houseYear := houseYearReg.FindString(houseInfo)
	houseFloorReg, _ := regexp.Compile(`共[1-9]\d*层`)
	houseFloor := houseFloorReg.FindString(houseInfo)
	return houseType, houseArea, houseOrientation, houseYear, houseFloor
}
