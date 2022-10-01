package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/zhulik/generic-tools/multiplexer"
	"github.com/zhulik/generic-tools/notification"
)

func main() {
	m := multiplexer.New[string]()
	defer m.Close()

	// for i := 0; i < 5; i++ {
	// 	go func(id int) {
	// 		for msg := range m.Subscribe() {
	// 			pp.Println("Receiver:", id, " msg: ", msg)
	// 		}
	// 	}(i)
	// }

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

		done := notification.New()

		go func() {
			pp.Println("Before receive")
			msg, ok := m.Receive()
			if !ok {
				panic("Not ok!")
			}
			pp.Println("Temporary receiver:", msg)
		}()

		m.Send(msg)
		done.Wait()
	}
}
