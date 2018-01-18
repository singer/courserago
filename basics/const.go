package main

const pi = 3.141
const (
	hello = "Привет"
	e     = 2.718
)
const (
	zero  = iota
	_     // пустая переменная, пропуск iota
	three // = 3
)
const (
	_         = iota             // пропускаем первое значне
	KB uint64 = 1 << (10 * iota) // 1024
	MB                           // 1048576
)
const (
	// нетипизированная константа
	year = 2017
	// типизированная константа
	yearTyped int = 2017
)

func main() {
	var month int32 = 13
	fmt.Println(month + year)

	// month + yearTyped (mismatched types int32 and int)
	// fmt.Println( month + yearTyped )
}
