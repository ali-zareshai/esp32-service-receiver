package domain

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"sensor_iot/Util"
	"strconv"
	"strings"
	"syscall"
)

func SetupMqtt() {
	topic := os.Getenv("MQTT_BROKER_TOPIC")
	opt := mqtt.NewClientOptions()
	opt.AddBroker(os.Getenv("MQTT_BROKER_URL"))
	opt.Username = os.Getenv("MQTT_BROKER_USERNAME")
	opt.Password = os.Getenv("MQTT_BROKER_PASSWORD")
	client := mqtt.NewClient(opt)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		Util.Logger.Error(token.Error().Error())
		return
	}

	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:", token.Error()))
	}

	fmt.Println("Subscribed to topic:", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Unsubscribe(topic)
	client.Disconnect(250)
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	msg := string(message.Payload())
	receiveMsgFromBroker(topic, msg)
}

func receiveMsgFromBroker(topic string, msg string) {
	var data DataModel
	parts := strings.Split(msg, ",")
	if len(parts) == 5 {
		data.Result = getParameter(parts[0])
		data.M32 = getParameter(parts[1])
		data.M33 = getParameter(parts[2])
		data.Ref = getParameter(parts[3])
		data.Vout = getParameter(parts[4])
		data.Device = strings.ReplaceAll(topic, "hello/", "")

		data.Save()
	} else {
		Util.Logger.Error("data is not complete")
	}
}

func getParameter(param string) float64 {
	num, err := strconv.ParseFloat(param, 64)
	if err != nil {
		Util.Logger.Error(err.Error())
		num = 0.0
	}
	return num
}
