package main

import (
	"gopkg.in/gomail.v2"
)

func main() {

	m := gomail.NewMessage()
	m.SetHeader("From", "xxxxxxxxxx@bbbb.com")
	m.SetHeader("To", "xxxxxxxxxx@bbbb.com")
	m.SetHeader("Subject", "Hello!")

	d := gomail.NewPlainDialer("smtp.exmail.qq.com", 465, "xxxxxxxxxx@bbbb.com", "xxxxxxx")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
