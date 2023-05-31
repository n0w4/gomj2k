package flags

import "flag"

type flags struct {
	BootstrapServer *string
	ToTopic         *string
	IsFreeMode      *bool
}

func Parse() *flags {
	flags := new(flags)

	flags.bootstrapServer()
	flags.isFreeMode()
	flags.toTopic()

	flag.Parse()

	return flags
}

func (f *flags) ValidateComposition() {
	if *f.IsFreeMode && *f.ToTopic == "" {
		panic("topic is required in free mode")
	}
}
