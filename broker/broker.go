package broker

import "github.com/n0w4/gomj2k/model"

type Broker interface {
	Publish(...model.StructuredMessage) error
}
