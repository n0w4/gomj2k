package serializations

import (
	"encoding/json"
	"fmt"

	"github.com/n0w4/gomj2k/model"
)

func RawToStructuredMessage(rawData []byte) (*model.StructuredMessage, error) {
	var bm model.StructuredMessage
	err := json.Unmarshal(rawData, &bm)
	if err != nil {
		return &bm, fmt.Errorf("error unmarshalling basic message: %v", err)
	}
	return &bm, nil
}
