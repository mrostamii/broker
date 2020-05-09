//
//  Hello World server.
//  Binds REP socket to tcp://*:5555
//  Expects "Hello" from client, replies with "World"
//

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	clintToServer, err := os.OpenFile("clientToServer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	serverToServerEnd, err := os.OpenFile("serverToServerEnd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind("tcp://*:5555")

	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://127.0.0.1:5556")

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		// _ = msg // to avoid declared but not used err.

		// write messages to a temp file
		_, err = clintToServer.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		_, err = clintToServer.Write([]byte("\n"))
		if err != nil {
			log.Fatal(err)
		}
		defer clintToServer.Close()

		//  Do some 'work'
		time.Sleep(time.Nanosecond)

		// send hello ------------------------------------------------
		msgEnd := fmt.Sprintf(msg)
		_, err = serverToServerEnd.Write([]byte(msgEnd))
		if err != nil {
			log.Fatal(err)
		}
		_, err = serverToServerEnd.Write([]byte("\n"))
		if err != nil {
			log.Fatal(err)
		}
		defer serverToServerEnd.Close()
		requester.Send(msgEnd, 0)

		// Wait for reply:
		replyEnd, _ := requester.Recv(0)
		_ = replyEnd
		// ------------------------------------------------

		//  Send reply back to client
		reply := "World"
		responder.Send(reply, 0)
		// fmt.Println("Sent ", reply)
	}
}
