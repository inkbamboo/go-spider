package consts

//go:generate gonum -types=HousePlatformEnum,SpiderTypeEnum
type HousePlatformEnum struct {
	Ke     string `enum:"ke,贝壳"`
	Anjuke string `enum:"anjuke,安居客"`
}
type SpiderTypeEnum struct {
	Area      string `enum:"area,区域"`
	ErShou    string `enum:"ershou,二手房"`
	ChengJiao string `enum:"chengjiao,二手成交"`
}
