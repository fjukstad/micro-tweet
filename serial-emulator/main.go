package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

type Message struct {
	From    string
	Message string
}

func main() {
	// /dev/cu.usbmodem1422
	c := &serial.Config{Name: "/dev/ptyp3", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	i := 0

	for {
		from := "From " + strconv.Itoa(i)
		message := "Message " + strconv.Itoa(i)
		m := Message{from, message}
		i = i + 1

		b, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = s.Write(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = s.Flush()
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(10 * time.Second)
	}

}
