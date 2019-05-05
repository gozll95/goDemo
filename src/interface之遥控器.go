package main

import (
	"fmt"
)

type Yaokongqi struct {
	Name string
}

type USB interface {
	selfName() string
	Connect()
	Disconnect()
}

func (Y Yaokongqi) selfName() string {
	return Y.Name
}

func (Y Yaokongqi) Connect() {
	fmt.Println("I will connect", Y.Name)
}

func (Y Yaokongqi) Disconnect() {
	fmt.Println("I will Disconnect", Y.Name)
}

func dis(usb USB) {
	if value, ok := usb.(Yaokongqi); ok {
		fmt.Println("dis", value.Name)
	} else {
		fmt.Println("Unknown device")
	}
}

func Connect(usb interface{}) {
	switch v := usb.(type) {
	case Yaokongqi:
		fmt.Println("Connect:", v.Name)
	default:
		fmt.Println("Unknow device")
	}
}

func main() {
	device := Yaokongqi{
		Name: "xxxxxxxxxx yaokongqi",
	}
	device.Connect()
	device.Disconnect()

	Connect(device)
	dis(device)
}
