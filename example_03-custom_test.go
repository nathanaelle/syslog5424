package syslog5424

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const TEST_SOCKET2 string = "./test-custom.socket"

func ExampleSyslogServerCustom() {
	defer os.Remove(TEST_SOCKET2)

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	mutex.Lock()

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	wg.Add(2)
	go server_custom(wg, mutex)
	go client_custom(wg, mutex)

	wg.Wait()
	mutex.Unlock()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing anoter stuff
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing a last stuff
}

func client_custom(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	sl_conn, err := (Dialer{
		QueueLen:   100,
		FlushDelay: 100 * time.Millisecond,
	}).Dial("unix", TEST_SOCKET2, new(T_RFC5425))
	if err != nil {
		log.Fatal(err)
	}

	syslog, err := New(sl_conn, LOG_DAEMON|LOG_WARNING, "custom-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	logger_err_conf := syslog.Channel(LOG_ERR).Logger("ERR : ")

	logger_err_conf.Print("doing some stuff")
	logger_err_conf.Print("doing anoter stuff")
	logger_err_conf.Print("doing a last stuff")

	sl_conn.End()
}

func server_custom(wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	collect, err := (Collector{
		QueueLen: 100,
	}).Collect("unix", TEST_SOCKET2, new(T_RFC5425))
	if err != nil {
		log.Fatal(err)
	}

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
