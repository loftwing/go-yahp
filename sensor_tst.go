package main

import (
	"github.com/loftwing/go-yahp/yahp/sensor"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	ports := []int{4444, 5555, 6666, 7777}
	sg := sensor.NewSensorGroup(ports...)
	sg.StartAll()

	go func(sg *sensor.SensorGroup) {
		mq := sg.Mq
		for {
			select {
			case msg := <-mq:
				log.Printf("message recvd: %+v\n", msg)
			default:
				time.Sleep(time.Second * 2)
			}
		}
	}(sg)

	wg.Add(1)
	wg.Wait()
}
