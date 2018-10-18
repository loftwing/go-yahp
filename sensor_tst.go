package main

import (
	"github.com/loftwing/go-yahp/yahp/sensor"
)

var (
	configServerAddr string
)

func CheckNetwork() {

}

func main() {
	forever := make(chan bool)



	lc := make(chan string, 512)
	ports := []int{4444, 5555, 6666, 7777}
	sg := sensor.NewPortManager(lc, ports...)
	sg.StartAll()

	<-forever
}
