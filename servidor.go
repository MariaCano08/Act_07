package main

import (
	//golang object
	"container/list"
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

func (t *Thread) run_to(n uint64) {
	for {
		if t.I == n || t.close {
			return
		}
		t.I++
		time.Sleep(time.Millisecond * 500)
	}

}

type Servidor struct {
	processes list.List
}

func (s *Servidor) add_process(t *Thread) {
	s.processes.PushBack(t)
}

func (s *Servidor) run() {
	const MAX_UNIT = ^uint64(0)
	const MAX_PROCESSES = 5

	for x := 1; x <= MAX_PROCESSES; x++ {
		t := Thread{ID: x, I: 0, close: false}
		go t.run_to(MAX_UNIT)
		s.add_process(&t)
	}

}

func (s *Servidor) GET() {
	server_GET, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		client, err := server_GET.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			s.response(client)
		}
	}
}

func (s *Servidor) POST() {
	const MAX_UNIT = ^uint64(0)
	server_POST, err := net.Listen("tcp", ":9998")

	if err != nil {
		fmt.Println(err)
		return
	}
	if server_POST != nil {
		for {
			client, err := server_POST.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				var p Thread
				err = gob.NewDecoder(client).Decode(&p)
				if err != nil {
					fmt.Println(err)
				} else {
					t := Thread{ID: p.ID, I: p.I, close: false}
					go t.run_to(MAX_UNIT)
					s.add_process(&t)
				}
			}
		}
	}
}

func (s *Servidor) response(client net.Conn) {
	process := s.processes.Front().Value
	//gob.Register(Thread{})
	//fmt.Println("enviando : ", process.(*Thread).i)
	//fmt.Println("enviando : ", process)

	err := gob.NewEncoder(client).Encode(process)
	if err != nil {
		fmt.Println(err)
		client.Close()
		return
	}

	client.Close()
	process.(*Thread).close = true
	s.processes.Remove(s.processes.Front())
}

func (s *Servidor) show() {

	for {
		for e := s.processes.Front(); e != nil; e = e.Next() {
			data := e.Value
			//fmt.Println(data)
			i := data.(*Thread).I
			id := data.(*Thread).ID
			fmt.Println(id, " : ", i)
		}
		time.Sleep(time.Second)
	}

}

func main() {

	server := Servidor{}
	server.processes.Init()

	server.run()
	go server.GET()
	go server.POST()
	go server.show()
	var input string
	fmt.Scanln(&input)
}
