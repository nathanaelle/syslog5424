package syslog5424

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const testTCPSocket string = "127.0.0.1:51400"

func ExampleTCPListener() {
	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	mutex.Lock()

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	wg.Add(2)
	go exTCPServer(wg, mutex)
	go exTCPClient(wg, mutex)

	wg.Wait()
	mutex.Unlock()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing anoter stuff
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing a last stuff
}

func exTCPClient(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	slConn, chanErr, err := Dial("tcp", testTCPSocket)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := <-chanErr; err != nil {
			log.Fatal(err)
		}
	}()

	syslog, err := New(slConn, LogDAEMON|LogWARNING, "client-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	loggerErrConf := syslog.Channel(LogERR).Logger("ERR : ")

	loggerErrConf.Print("doing some stuff")
	loggerErrConf.Print("doing anoter stuff")
	loggerErrConf.Print("doing a last stuff")

	slConn.End()
}

func exTCPServer(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	listener, err := TCPListener("tcp", testTCPSocket)
	if err != nil {
		log.Fatal(err)
	}

	collect, chanErr := NewReceiver(listener, 100, TransportLFEnded)

	go func() {
		if err := <-chanErr; err != nil {
			log.Fatalf("client chanErr %q", err)
		}
	}()

	// socket is created
	mutex.Unlock()

	count := 3
	for count > 0 {
		count--

		msg, err, _ := collect.Receive()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", msg.String())
	}

	collect.End()
}
