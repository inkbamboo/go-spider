package util

import (
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strings"
)

func GetHouseType(houseInfo string) string {
	houseTypeReg, _ := regexp.Compile(`\d+室\d+厅`)
	houseType := houseTypeReg.FindString(houseInfo)
	return houseType
}
func GetHouseArea(houseInfo string) float64 {
	houseAreaReg, _ := regexp.Compile(`[1-9]\d*(\.\d*)?平米?`)
	houseArea := cast.ToFloat64(strings.Replace(houseAreaReg.FindString(houseInfo), "平米", "", -1))
	return houseArea
}
func GetHouseOrientation(houseInfo string) string {
	houseOrientationReg, _ := regexp.Compile(`[东北西南]`)
	houseOrientation := strings.Join(houseOrientationReg.FindAllString(houseInfo, -1), "")
	return houseOrientation
}
func GetHouseYear(houseInfo string) string {
	houseYearReg, _ := regexp.Compile(`[1-9]\d*年`)
	houseYear := houseYearReg.FindString(houseInfo)
	return houseYear
}
func GetHouseFloor(houseInfo string) string {
	houseFloorReg, _ := regexp.Compile(`共[1-9]\d*层`)
	houseFloor := houseFloorReg.FindString(houseInfo)
	return houseFloor
}
func TrimInfoEmpty(houseInfo string) string {
	emptyReg, _ := regexp.Compile(`\s`)
	houseInfo = emptyReg.ReplaceAllString(houseInfo, "")
	return houseInfo
}
func GetTotalPrice(houseInfo string) float64 {
	totalPriceReg, _ := regexp.Compile(`[1-9]\d*(\.\d*)?万?`)
	housePrice := totalPriceReg.FindString(houseInfo)
	fmt.Printf("housePrice: %s\n", housePrice)
	return cast.ToFloat64(strings.Replace(housePrice, "万", "", -1))
}

func GetUnitPrice(houseInfo string) float64 {
	unitPrice := strings.ReplaceAll(houseInfo, "元/平", "")
	unitPrice = strings.TrimSpace(strings.ReplaceAll(unitPrice, ",", ""))
	return cast.ToFloat64(unitPrice)
}
func GetDealCycle(houseInfo string) int64 {
	houseFloorReg, _ := regexp.Compile(`\d*天?`)
	dealCycle := houseFloorReg.FindString(houseInfo)
	return cast.ToInt64(strings.Replace(dealCycle, "天", "", -1))
}
