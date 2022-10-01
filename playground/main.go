package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/zhulik/generic-tools/multiplexer"
)

func main() {
	m := multiplexer.New[string]()
	defer m.Close()

	for i := 0; i < 5; i++ {
		go func(id int) {
			for msg := range m.Subscribe() {
				pp.Println("Receiver:", id, "msg:", msg)
			}
		}(i)
	}

	var reader = bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		msg = strings.TrimSuffix(msg, "\n")
		if msg == "exit" {
			break
		}

		m.Send(msg)
	}
}
