package consts

//go:generate gonum -types=PoetryPlatformEnum,PoetryTypeEnum,SpiderTypeEnum
type PoetryPlatformEnum struct {
	ZHSC string `enum:"zhsc,中华诗词"`
}
type PoetryTypeEnum struct {
	Shi string `enum:"shi,诗"`
	Ci  string `enum:"ci,词"`
	Wen string `enum:"wen,文"`
	Qu  string `enum:"qu,曲"`
	Fu  string `enum:"fu,赋"`
}
type SpiderTypeEnum struct {
	Author string `enum:"author,作者"`
	Poetry string `enum:"poetry,诗词"`
}
