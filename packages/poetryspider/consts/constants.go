package consts

//go:generate gonum -types=PoetryPlatformEnum
type PoetryPlatformEnum struct {
	ZHSC string `enum:"zhsc,中华诗词"`
}
