package yahp

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

type Sensor struct {
	Port    int
	Started bool
}

type SensorGroup struct {
	Mq      chan string
	Sensors []*Sensor
}

func NewSensorGroup(port ...int) *SensorGroup {
	var sensors []*Sensor
	for _, v := range port {
		sensor := &Sensor{Port: v, Started: false}
		sensors = append(sensors, sensor)
	}
	return &SensorGroup{
		Mq:      make(chan string, 512),
		Sensors: sensors}
}

func (sg *SensorGroup) StartAll() {
	for _, v := range sg.Sensors {
		go v.start(sg.Mq)
	}
}

func startListener(ctx context.Context, port string) (string, error) {
	cmd := exec.CommandContext(ctx, "port.exe", port)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (s *Sensor) start(mc chan<- string) {
	for {
		s.Started = true
		port := fmt.Sprintf("%d", s.Port)
		out, err := startListener(context.Background(), port)
		if err != nil {
			log.Println("sensor exec err: ", err)
		}

		// set started to false, send output through to channel of sgroup
		s.Started = false
		mc <- out

		// wait 5 mins before re-opening the port, disabled while developing
		//time.Sleep(time.Minute * 5)
	}
}
