package flags

import "flag"

var freeModeUsage string = `if true, you will need use "to-topic" option, 
and I send from stdin as is. (default: false)`

func (f *flags) isFreeMode() {
	var ifm bool
	flag.BoolVar(&ifm, "free-mode", false, freeModeUsage)
	f.IsFreeMode = &ifm
}
