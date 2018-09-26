package sensor

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Port struct {
	Port    int
	Started bool
}

type PortManager struct {
	Mq      chan string
	LogChan chan<- string
	Ports   []*Port
}

func NewPortManager(lc chan<- string, port ...int) *PortManager {
	var ports []*Port
	for _, v := range port {
		sensor := &Port{Port: v, Started: false}
		ports = append(ports, sensor)
	}
	return &PortManager{
		Mq:      make(chan string, 512),
		Ports:   ports,
		LogChan: lc,
	}
}

func (pm *PortManager) StartAll() {
	for _, v := range pm.Ports {
		go v.start(pm.Mq)
	}
	pm.StartMonitor()
}

func startListener(ctx context.Context, port string) (string, error) {
	cmd := exec.CommandContext(ctx, "port.exe", port)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (pm *PortManager) StartMonitor() {
	go func(lc chan<- string) {
		for {
			for _, v := range pm.Ports {
				lc <- fmt.Sprintf("%d:%t", v.Port, v.Started)
			}
			time.Sleep(time.Second * 150)
		}
	}(pm.LogChan)
}

func (s *Port) start(mc chan<- string) {
	for {
		s.Started = true
		port := fmt.Sprintf("%d", s.Port)
		out, err := startListener(context.Background(), port)
		if err != nil {
			log.Println("sensor exec err: ", err)
			s.Started = false
		}

		// send output through to channel of sgroup
		mc <- out

		// wait 1 min before re-opening the port, disabled while developing
		time.Sleep(time.Minute)
	}
}
