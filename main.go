package main

import (
	`bufio`
	`context`
	`fmt`
	`gopkg.in/natefinch/npipe.v2`
	`io`
	`log`
	`net`
	`os/exec`
	`sync`
	`time`
)

)

var wg sync.WaitGroup

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func startListener(ctx context.Context, port string) {
	cmd := exec.CommandContext(ctx, "port.exe", port)
	err := cmd.Run()
	if err != nil {
		failOnError(err, "Run listener")
	}
}

func pipeHandler(c net.Conn) {
	r := bufio.NewReader(c.)
	for {
		log.Printf("Buffered bytes: %d\n", r.Buffered())
		if r.Buffered() > 1 {
			var read = make([]byte, r.Buffered())
			cRead, err := r.Read(read)
			if err != io.EOF {
				log.Println("READ ", cRead, " bytes: ", string(read))
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func startPipeServer() {
	ln, err := npipe.Listen(`\\.\pipe\yahp`)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()
		log.Printf("Pipe connection created: %+v\n", conn)
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
}

func main() {
	go startPipeServer()

	wg.Add(1)
	//go respawn("5555")
	//go respawn("6666")
	//go respawn("7777")
	//go respawn("8888")
	wg.Wait()

	log.Println("All finished? uhh")
}
