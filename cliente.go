package main

import (
	//golang object

	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Thread struct {
	ID    int
	I     uint64
	close bool
}

func (c *Cliente) run_to(n uint64) {
	for {
		if c.process.I == n || c.process.close {
			return
		}
		c.process.I++
		time.Sleep(time.Millisecond * 500)
	}

}

type Cliente struct {
	process Thread
}

func (c *Cliente) GET() {
	const MAX_UNIT = ^uint64(0)

	cliente, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println("e1 : ", err)
		return
	}

	//var count uint64
	//err = gob.NewDecoder(client).Decode(&count)

	var p Thread
	err = gob.NewDecoder(cliente).Decode(&p)

	if err != nil {
		fmt.Println("e2 : ", err)
	} else {
		c.process = Thread{ID: p.ID, I: p.I, close: false}
		go c.run_to(MAX_UNIT)
		go c.show()
	}
}

func (c *Cliente) POST() {
	client, err := net.Dial("tcp", ":9998")

	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(client).Encode(&c.process)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Cliente) show() {

	for {
		i := c.process.I
		id := c.process.ID
		fmt.Println(id, " : ", i)
		time.Sleep(time.Second)
	}

}

func main() {
	cli := Cliente{}
	go cli.GET()
	var input string
	fmt.Scanln(&input)
	cli.POST()
	cli.process.close = true
}
