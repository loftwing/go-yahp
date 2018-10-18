package sensor

import (
	`log`

	`golang.org/x/sys/windows/registry`
)

type SensorManager struct {
	ConfigServerAddr string

	LogChan  chan string

	Pm *PortManager
	Mm *MessageManager
}

func NewSensorManager() (*SensorManager) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\yahp`, registry.QUERY_VALUE)
	if err != nil {
		log.Panic("REGISTRY: yahp key could not be opened.")
	}
	defer k.Close()

	s, _, err := k.GetStringValue("ConfigServerAddr")
	if err != nil {
		log.Panic("REGISTRY: Value ConfigServerAddr could not be read.")
	}

	return &SensorManager{
		ConfigServerAddr: s,
		Pm: nil,
		Mm: nil,
	}
}

func getConfig() {

}

func (sm *SensorManager) Run() {

}
