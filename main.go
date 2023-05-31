package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/n0w4/gomj2k/broker"
	"github.com/n0w4/gomj2k/flags"
	"github.com/n0w4/gomj2k/model"
	"github.com/n0w4/gomj2k/serializations"
)

func main() {

	fs := flags.Parse()
	fs.ValidateComposition()

	myBroker := broker.NewKafka(*fs.BootstrapServer).WithProducer()

	fromStdin(fs.ToTopic, fs.IsFreeMode, myBroker)

}

func fromStdin(givenTopic *string, isFreeMode *bool, bk broker.Broker) {
	reader := bufio.NewReader(os.Stdin)

	for {
		stdin, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(">>>> end <<<<")
				break
			}
			fmt.Printf("error reading from stdin: %v\n", err)
			os.Exit(1)
		}

		line, err := consolidateLine(stdin)
		if err != nil {
			log.Printf("error consolidating line: %v", err)
		}

		messages, err := freeOrStructured(isFreeMode, line, givenTopic)
		if err != nil {
			log.Printf("whe serialize %v", err)
			continue
		}

		bk.Publish(messages...)
	}
}

func consolidateLine(rawLine []byte) ([]byte, error) {
	line := string(rawLine)

	line = removeBreakLine(line)

	if isEmptyLine(line) {
		return nil, fmt.Errorf("empty line")
	}
	return []byte(line), nil
}

func removeBreakLine(line string) string {
	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	return line
}

func isEmptyLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

func freeOrStructured(isFreeMode *bool, line []byte, topic *string) ([]model.StructuredMessage, error) {
	messages := make([]model.StructuredMessage, 0)

	if *isFreeMode {
		return freeMode(line, topic)
	}

	structuredMessage, err := serializations.RawToStructuredMessage(line)
	if err != nil {
		return nil, err
	}
	messages = append(messages, *structuredMessage)

	if topic != nil && *topic != "" {
		structuredMessage, err := serializations.RawToStructuredMessage(line)
		if err != nil {
			return nil, err
		}
		structuredMessage.Topic = *topic
		messages = append(messages, *structuredMessage)
	}

	return messages, nil
}

func freeMode(line []byte, topic *string) ([]model.StructuredMessage, error) {
	messages := make([]model.StructuredMessage, 0)
	if topic != nil {
		message := model.StructuredMessage{
			Payload: line,
			Topic:   *topic,
		}
		messages = append(messages, message)
		return messages, nil
	}
	return nil, fmt.Errorf("need a topic")
}
