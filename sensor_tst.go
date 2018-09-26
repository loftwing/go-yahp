package main

import (
	"github.com/loftwing/go-yahp/yahp/sensor"
	"log"
)

func main() {
	forever := make(chan bool)

	ports := []int{4444, 5555, 6666, 7777}
	sg := sensor.NewSensorGroup(ports...)
	sg.StartAll()

	mm, err := sensor.NewMessageManager(
		"amqp://guest:guest@localhost:5672/",
		"yahp",
		sg.Mq,
	)
	if err != nil {
		log.Panic(err)
	}
	mm.Start()

	<-forever
}
