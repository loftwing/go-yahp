package yahp

import (
	`encoding/json`
	`strings`
)

type MessageManager struct {
	sg       *SensorGroup

}

type SensorMessage struct {
	Returntype string `json:"returntype"`
	Port       int64  `json:"port"`
	Ip         string `json:"ip"`
}

func (mm *MessageManager) Start() {
	go func (mq <-chan string) {
		
	} (mm.sg.mc)
}

func validateMessage(msg *SensorMessage) bool {
	return (len(msg.Ip) <= 16) &&
		(msg.Port > 1 && msg.Port < 65535) &&
		(msg.Returntype == "con" || msg.Returntype == "err")
}

func ParseSensorMessage(msg string) (SensorMessage, error) {
	dec := json.NewDecoder(strings.NewReader(msg))
	sm := SensorMessage{}
	err := dec.Decode(sm)
	if err != nil {
		return SensorMessage{}, err
	} else {
		return sm, nil
	}
}
