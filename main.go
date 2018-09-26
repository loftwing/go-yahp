package main

import (
	"log"
	"sync"
	"time"

	"github.com/loftwing/go-yahp/yahp"
)

var wg sync.WaitGroup

func main() {
	ports := []int{4444, 5555, 6666, 7777}
	sg := yahp.NewSensorGroup(ports...)
	sg.StartAll()

	go func(sg *yahp.SensorGroup) {
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
