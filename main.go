package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"sync"

	"gopkg.in/natefinch/npipe.v2"
)

// var elog debug.Log
//
// type MinisocService struct {}
//
// func (s *MinisocService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
// 	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
//
// 	changes <- svc.Status{State: svc.StartPending}
//
// 	fasttick := time.Tick(500 * time.Millisecond)
// 	slowtick := time.Tick(2 * time.Second)
// 	tick := fasttick
//
// 	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
//
// 	elog.Info(1, strings.Join(args, "-"))
//
// 	loop:
// 		for {
// 			select {
// 			case <-tick:
// 				elog.Info(1, "pull stdout")
// 			case c := <-r:
// 				switch c.Cmd {
// 				case svc.Interrogate:
// 					changes <- c.CurrentStatus
// 					time.Sleep(100 * time.Millisecond)
// 					changes <- c.CurrentStatus
// 				case svc.Stop, svc.Shutdown:
// 					break loop
// 				case svc.Pause:
// 					changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
// 					tick = slowtick
// 				case svc.Continue:
// 					changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
// 					tick = fasttick
// 				default:
// 					elog.Error(1, fmt.Sprintf("weird control request #%d", c))
// 				}
// 			}
// 		}
// 	changes <- svc.Status{State: svc.StopPending}
// 	return
// }
//
// func runService (name string, isDebug bool) {
// 	var err error
// 	if isDebug {
// 		elog = debug.New(name)
// 	} else {
// 		elog, err = eventlog.Open(name)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	defer elog.Close()
//
// 	elog.Info(1, fmt.Sprintf("starting %s", name))
//
// 	run := svc.Run
//
// 	if isDebug {
// 		run = debug.Run
// 	}
//
// 	err = run(name, &MinisocService{})
// 	if err != nil {
// 		elog.Error(1, fmt.Sprintf("%s stopped", name))
// 	}
// }

var wg sync.WaitGroup

type SensorMessage struct {
	Returntype string `json:"returntype"`
	Port       int64  `json:"port"`
	Ip         string `json:"ip"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func startListener(ctx context.Context, port string) {
	cmd := exec.CommandContext(ctx, "port.exe", port)
	err2 := cmd.Start()
	if err2 != nil {
		log.Println("Failed exec cmd")
	}
}

func pipeHandler(c net.Conn) {
	//r := bufio.NewReader(c)

}

func startPipeServer() {
	ln, err := npipe.Listen(`\\.\pipe\yahp`)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[ERR] PIPE: ", err)
			continue
		}

		go pipeHandler(conn)
	}
}

func respawn(port string) {
	for {
		startListener(context.Background(), port)
	}
	wg.Done()
}

func main() {
	wg.Add(4)
	go respawn("5555")
	go respawn("6666")
	go respawn("7777")
	go respawn("8888")
	wg.Wait()
	log.Println("All finished? uhh")
}