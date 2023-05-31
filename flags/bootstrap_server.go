package flags

import "flag"

var bootstrapServerUsage string = `bootstrap server to connect to.`

func (f *flags) bootstrapServer() {
	var bs string
	flag.StringVar(&bs, "bs", "localhost:9092", bootstrapServerUsage)
	f.BootstrapServer = &bs
}
