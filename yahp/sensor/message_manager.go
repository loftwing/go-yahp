package sensor

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"strings"
)

type MessageManager struct {
	RmqConnection *amqp.Connection
	RmqChannel    *amqp.Channel
	RmqConnString string
	RmqExchange   string
	SensorMq      <-chan string
}

type SensorMessage struct {
	Returntype string `json:"returntype"`
	Port       int64  `json:"port"`
	Ip         string `json:"ip"`
}

func NewMessageManager(connstring, exchange string, mq <-chan string) (*MessageManager, error) {
	mm := &MessageManager{
		RmqConnection: nil,
		RmqChannel:    nil,
		RmqConnString: connstring,
		RmqExchange:   exchange,
		SensorMq:      mq,
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
			if sm, err := ParseSensorMessage(msg); err == nil {
				if valid := validateMessage(sm); valid == true {
					err := mm.Publish(sm, "yahp.connection")
					if err != nil {
						log.Println("[MM] Failed to publish message to rmq: ", err)
					}
				} else {
					log.Printf("[MM] Sensor message failed to validate: %+v\n", sm)
				}
			} else {
				log.Println("[MM] Failed to parse message from sensor: ", msg, " error: ", err)
			}
		}
	}(mm.SensorMq)
}

func (mm *MessageManager) Publish(msg *SensorMessage, topic string) error {
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

func validateMessage(msg *SensorMessage) bool {
	return (len(msg.Ip) <= 16) &&
		(msg.Port > 1 && msg.Port < 65535) &&
		(msg.Returntype == "con" || msg.Returntype == "err")
}

func ParseSensorMessage(msg string) (*SensorMessage, error) {
	dec := json.NewDecoder(strings.NewReader(msg))
	sm := &SensorMessage{}
	err := dec.Decode(sm)
	if err != nil {
		return nil, err
	}

	if !validateMessage(sm) {
		return sm, errors.New("message validation failed")
	}

	return sm, nil
}
