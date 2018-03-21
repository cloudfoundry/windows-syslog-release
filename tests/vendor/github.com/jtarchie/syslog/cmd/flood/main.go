package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
	"sync/atomic"
)

var failedMessages int32

func main() {
	workers := flag.Int("num-workers", 1, "number of workers to generate messages")
	numMessages := flag.Int("num-messages", 1000, "number of messages to send")
	uri := flag.String("uri", "", "the endpoint to connect to (udp://)")
	flag.Parse()

	wait := sync.WaitGroup{}

	for i := 0; i < *workers; i++ {
		wait.Add(1)
		go func(worker int) {
			defer wait.Done()
			setupWorker(worker, numMessages, uri)
		}(i)
	}

	wait.Wait()

	fmt.Printf("failed sending %d messages", failedMessages	)
}

func setupWorker(worker int, numMessages *int, uri *string) {
	fmt.Printf("Starting worker %d with %d messages\n", worker, *numMessages)
	conn := createConnection(uri)
	defer conn.Close()
	for j := 0; j < *numMessages; j++ {
		err := sendMessage(conn, worker, j)
		if err != nil {
			atomic.AddInt32(&failedMessages, 1)
		}
	}
}

func sendMessage(conn net.Conn, worker int, j int) error {
	_, err := fmt.Fprintf(
		conn,
		"<%d>%d %s worker-%d.example.com %s %d %d - This is message %d from worker %d\n",
		rand.Int31n(192), //priority 0-191
		1,                //version
		time.Now().Format("2006-01-01T15:04:05.000Z"),
		worker,      //host
		"flood",     //app name
		os.Getpid(), //proc id
		j,           //log id
		j, worker,
	)
	return err
}
func createConnection(uri *string) net.Conn {
	conn, err := net.Dial("udp", *uri)
	if err != nil {
		log.Fatalf("cannot connect: %s", err)
	}
	return conn
}
