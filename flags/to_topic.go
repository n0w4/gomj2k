package flags

import "flag"

var toTopicUsage string = `what topic do you want to send to. 
(required if not present on stdin)`

func (f *flags) toTopic() {
	var tt string
	flag.StringVar(&tt, "to-topic", "", toTopicUsage)
	if tt == "" {
		f.ToTopic = nil
	}
	f.ToTopic = &tt
}
