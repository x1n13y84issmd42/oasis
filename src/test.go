package main

import (
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func main() {
	logger := log.New("festive", 6)
	logger.NOMESSAGE("Hello")

	blogger1 := log.NewBuffer(log.New("festive", 6))
	blogger2 := log.NewBuffer(log.New("festive", 6))
	blogger3 := log.NewBuffer(log.New("festive", 6))
	blogger4 := log.NewBuffer(log.New("festive", 6))

	blogger1.NOMESSAGE("Hello 1")
	blogger3.NOMESSAGE("Hello 3")
	blogger2.NOMESSAGE("Hello 2")

	blogger2.NOMESSAGE("World 2")
	blogger1.NOMESSAGE("World 1")
	blogger3.NOMESSAGE("World 3")
	blogger1.NOMESSAGE("YOLO 1")
	blogger3.NOMESSAGE("YOLO 3")

	blogger4.NOMESSAGE("Hello 4")
	blogger4.NOMESSAGE("World 4")
	blogger2.NOMESSAGE("YOLO 2")
	blogger4.NOMESSAGE("YOLO 4")

	blogger1.Flush()
	blogger2.Flush()
	blogger3.Flush()
	blogger4.Flush()
}
