package syslog5424

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const TEST_SOCKET string = "./test.socket"

func ExampleSyslogServer() {
	defer os.Remove(TEST_SOCKET)

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	mutex.Lock()

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	wg.Add(2)
	go server(wg, mutex)
	go client(wg, mutex)

	wg.Wait()
	mutex.Unlock()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing anoter stuff
	// <27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing a last stuff
}

func client(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	slConn, chanErr, err := Dial("unix", TEST_SOCKET)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := <-chanErr; err != nil {
			log.Fatal(err)
		}
	}()

	syslog, err := New(slConn, LOG_DAEMON|LOG_WARNING, "client-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	logger_err_conf := syslog.Channel(LOG_ERR).Logger("ERR : ")

	logger_err_conf.Print("doing some stuff")
	logger_err_conf.Print("doing anoter stuff")
	logger_err_conf.Print("doing a last stuff")

	slConn.End()
}

func server(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	listener, err := UnixListener(TEST_SOCKET)
	if err != nil {
		log.Fatal(err)
	}
	collect, chanErr := NewReceiver(listener, 100, T_ZEROENDED)

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
