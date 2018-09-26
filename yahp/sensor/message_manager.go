package sensor

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"net"
	"strings"
)

type MessageManager struct {
	RmqConnection *amqp.Connection
	RmqChannel    *amqp.Channel
	RmqConnString string
	RmqExchange   string
	SensorMq      <-chan string
	LogChan       <-chan string
	HostIP        string
}

type PortMessage struct {
	Returntype string `json:"returntype"`
	Port       int64  `json:"port"`
	Ip         string `json:"ip"`
}

type Message struct {
	HostIP string
	Pm     *PortMessage
}

func NewMessageManager(connstring, exchange string, mq, lc <-chan string) (*MessageManager, error) {
	mm := &MessageManager{
		RmqConnection: nil,
		RmqChannel:    nil,
		RmqConnString: connstring,
		RmqExchange:   exchange,
		SensorMq:      mq,
		LogChan:       lc,
		HostIP:        localAddress(),
	}

	err := mm.Connect()
	if err != nil {
		return nil, err
	}

	if mm.RmqChannel != nil && mm.RmqConnection != nil {
		return mm, nil
	} else {
		return nil, errors.New("Failed to create mm. Unknown error.")
	}
}

func (mm *MessageManager) Connect() error {
	var err error = nil
	// Connect to rmq broker and set mm pointer
	mm.RmqConnection, err = amqp.Dial(mm.RmqConnString)
	if err != nil {
		return err
	}
	log.Println("[RMQ] Connected to RabbitMQ")

	// Create channel and set mm pointer
	mm.RmqChannel, err = mm.RmqConnection.Channel()
	if err != nil {
		return err
	}
	log.Println("[RMQ] Channel created")

	// Declare the exchange using mm pointer
	err = mm.RmqChannel.ExchangeDeclare(
		mm.RmqExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mm *MessageManager) Start() {
	// goroutine that constantly watches the channel for output from sensors
	go func(mq <-chan string) {
		for {
			msg := <-mq
			if m, err := mm.ParseSensorMessage(msg); err == nil {
				if valid := validateMessage(m.Pm); valid == true {
					err := mm.Publish(m, "yahp.connection")
					if err != nil {
						log.Println("[MM] Failed to publish message to rmq: ", err)
					}
				} else {
					log.Printf("[MM] Port message failed to validate: %+v\n", m)
				}
			} else {
				log.Println("[MM] Failed to parse message from sensor: ", msg, " error: ", err)
			}
		}
	}(mm.SensorMq)

	go func(lc <-chan string) {
		for {
			log := <-lc
			mm.PublishLog(log)
		}
	}(mm.LogChan)
}

func (mm *MessageManager) PublishLog(msg string) {
	// TODO log msg stored in ip field for now...
	ts := "yahp.log.listen"
	sm := &PortMessage{
		Returntype: "log",
		Port:       1,
		Ip:         msg,
	}

	m := &Message{HostIP: mm.HostIP, Pm: sm}
	if err := mm.Publish(m, ts); err != nil {
		log.Println("[LOG] Failed to publish log to rmq!")
	}
}

func (mm *MessageManager) Publish(msg *Message, topic string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = mm.RmqChannel.Publish(
		mm.RmqExchange,
		topic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func validateMessage(msg *PortMessage) bool {
	return (len(msg.Ip) <= 16) &&
		(msg.Port > 1 && msg.Port < 65535) &&
		(msg.Returntype == "con" || msg.Returntype == "err")
}

func (mm *MessageManager) ParseSensorMessage(msg string) (*Message, error) {
	dec := json.NewDecoder(strings.NewReader(msg))
	sm := &PortMessage{}
	err := dec.Decode(sm)
	if err != nil {
		return nil, err
	}

	if !validateMessage(sm) {
		return nil, errors.New("message validation failed")
	}

	return &Message{HostIP: mm.HostIP, Pm: sm}, nil
}

func localAddress() string {
	ip := "x.x.x.x"
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("[ERR] Cant get interface address")
		return ip
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, a := range addrs {
			addr := a.String()
			if !strings.Contains(addr, ":") && !strings.Contains(addr, "127.0.0.1") {
				return addr
			}
		}
	}
	return ip
}
