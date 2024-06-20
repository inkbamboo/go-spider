package util

import (
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

func ParseHouseDetail(houseInfo string) (float64, string, string, string, string) {
	emptyReg, _ := regexp.Compile(`\s`)
	houseInfo = emptyReg.ReplaceAllString(houseInfo, "")
	houseTypeReg, _ := regexp.Compile(`\d+室\d+厅`)
	houseType := houseTypeReg.FindString(houseInfo)
	houseAreaReg, _ := regexp.Compile(`[1-9]\d*(\.\d*)?平米`)
	houseArea := cast.ToFloat64(strings.Replace(houseAreaReg.FindString(houseInfo), "平米", "", -1))
	houseOrientationReg, _ := regexp.Compile(`[东北西南]`)
	houseOrientation := strings.Join(houseOrientationReg.FindAllString(houseInfo, -1), "")
	houseYearReg, _ := regexp.Compile(`[1-9]\d*年`)
	houseYear := houseYearReg.FindString(houseInfo)
	houseFloorReg, _ := regexp.Compile(`共[1-9]\d*层`)
	houseFloor := houseFloorReg.FindString(houseInfo)
	return houseArea, houseType, houseOrientation, houseYear, houseFloor
}
