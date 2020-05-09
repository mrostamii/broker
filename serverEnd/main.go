//
//  Hello World server.
//  Binds REP socket to tcp://*:5555
//  Expects "Hello" from client, replies with "World"
//

package main

import (
	"log"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	clintToEndServer, err := os.OpenFile("clientToEndServer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind("tcp://*:23001")

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		// _ = msg // to avoid declared but not used err.

		// write messages to a temp file
		_, err = clintToEndServer.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		_, err = clintToEndServer.Write([]byte("\n"))
		if err != nil {
			log.Fatal(err)
		}
		defer clintToEndServer.Close()
		// fmt.Println(msg)

		//  Do some 'work'
		time.Sleep(time.Nanosecond)

		//  Send reply back to client
		reply := "World"
		responder.Send(reply, 0)
		// fmt.Println("Sent ", reply)
	}
}
