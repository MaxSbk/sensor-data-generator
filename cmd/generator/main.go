package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sensor-data-generator/internal/config"
	"sensor-data-generator/internal/mqtt"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	conf, err := config.LoadConfig("./configs/application.yaml")
	if err != nil {
		panic(err)
	}

	uri, err := url.Parse(conf.Mqtt.Url)
	if err != nil {
		panic(err)
	}

	client := mqtt.CreateClient(conf.Mqtt.ClientId, uri)
	defer func() {
		fmt.Println("Closing MQTT client...")
		client.Disconnect(100)
	}()

	for _, sensor := range conf.Sensors {
		go generate(client, conf.Mqtt.Topic, conf.Mqtt.Qos, sensor)
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}

func generate(client paho.Client, topic string, qos byte, sensor config.Sensor) {
	var sensorRow = "type:%s,id:%s,machine-id:%s,part-id:%s,tool-id:%s,timestamp:%v,unit:%s,value:%v"

	interval := sensor.Generator.Interval
	timer := time.NewTicker(time.Duration(interval) * time.Millisecond)
	min := sensor.Generator.Range.Min
	max := sensor.Generator.Range.Max

	var counter uint
	var extraFlag byte = 0     // 0 - normal, 1 - below, 2 - above
	var lastExtraFlag byte = 0 // 0 - normal, 1 - below, 2 - above
	for range timer.C {
		sensorValue := min + rand.Float64()*(max-min)

		counter += 1

		sensorValue, counter, extraFlag, lastExtraFlag = extraValueLogic(sensorValue, counter, extraFlag, lastExtraFlag, sensor.Generator)

		payload := fmt.Sprintf(sensorRow, sensor.Type, sensor.Id, sensor.MachineId, sensor.PartId, sensor.ToolId, time.Now().Unix(), sensor.Unit, sensorValue)
		fmt.Println(payload)
		client.Publish(topic, qos, false, payload)
	}
}

func extraValueLogic(sensorValue float64, counter uint, extraFlag byte, lastExtraFlag byte, config config.Generator) (float64, uint, byte, byte) {
	if extraFlag == 0 {
		if config.ExtraBelowValues.Freq != 0 && lastExtraFlag != 1 && counter%config.ExtraBelowValues.Freq == 0 {
			extraFlag = 1
			counter = 1
		} else if config.ExtraAboveValues.Freq != 0 && lastExtraFlag != 2 && counter%config.ExtraAboveValues.Freq == 0 {
			extraFlag = 2
			counter = 1
		}
	}

	if extraFlag == 1 {
		if counter == config.ExtraBelowValues.Duration {
			extraFlag = 0
			lastExtraFlag = 1
			counter = 1
		}
		sensorValue = (config.Range.Min - config.Range.Min*config.ExtraBelowValues.PercentageDeviation/100) + rand.Float64()
	} else if extraFlag == 2 {
		if counter == config.ExtraAboveValues.Duration {
			extraFlag = 0
			lastExtraFlag = 2
			counter = 1
		}
		sensorValue = (config.Range.Max + config.Range.Max*config.ExtraAboveValues.PercentageDeviation/100) + rand.Float64()
	}
	return sensorValue, counter, extraFlag, lastExtraFlag
}
