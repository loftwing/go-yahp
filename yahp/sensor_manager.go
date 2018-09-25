package yahp

import (
	`context`
	`log`
	`os/exec`
	`time`
)

type Sensor struct {
	Port     int
	Started  bool
}

type SensorGroup struct {
	mc       chan string
	Sensors  []*Sensor
}

func NewSensorGroup(port ...int) (*SensorGroup) {
	var sensors []*Sensor
	for _, v := range port {
		sensor := &Sensor{Port: v, Started: false}
		sensors = append(sensors, sensor)
	}
	return &SensorGroup{
		mc: make(chan string, 512),
		Sensors: sensors}
}

func (sg *SensorGroup) StartAll() {
	for _,v := range sg.Sensors {
		go v.start(sg.mc)
	}
}

func startListener(ctx context.Context, port string) (string, error) {
	cmd := exec.CommandContext(ctx, "port.exe", port)
	out, err := cmd.Output()
	if err != nil {
		log.Println("Failed exec cmd")
		return "", err
	}

	return string(out), nil
}

func (s *Sensor) start(mc chan<- string) {
	for {
		s.Started = true
		out, err := startListener(context.Background(), string(s.Port))
		if err != nil {
			log.Println("error returning from sensor: ", err)
		}

		// set started to false, send output through to channel of sgroup
		s.Started = false
		mc <- out

		// wait 5 mins before re-opening the port
		time.Sleep(time.Minute * 5)
	}
}
