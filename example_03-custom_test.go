package syslog5424_test

import (
	//"."
	"os"
	"log"
	"fmt"
	"time"
	"sync"
	"github.com/nathanaelle/syslog5424"
)


const	TEST_SOCKET2 string = "./test-custom.socket"



func ExampleSyslogServerCustom() {
	defer os.Remove(TEST_SOCKET2)

	wg	:= new(sync.WaitGroup)
	mutex	:= new(sync.Mutex)

	mutex.Lock()

	syslog5424.Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}


	go server_custom(wg, mutex)
	go client_custom(wg, mutex)

	time.Sleep(time.Second)
	wg.Wait()
	mutex.Unlock()

	// Output:
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing some stuff
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing anoter stuff
	// <27>1 2014-12-20T14:04:00Z localhost custom-app 1234 - - ERR : doing a last stuff
}


func client_custom(wg *sync.WaitGroup, mutex *sync.Mutex)  {
	wg.Add(1)
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	sl_conn,err := (syslog5424.Dialer{
		QueueLen:	100,
		FlushDelay:	100*time.Millisecond,
	}).Dial("unix", TEST_SOCKET2, new(syslog5424.T_RFC5426) )
	if err != nil {
		log.Fatal(err)
	}

	syslog, err := syslog5424.New(sl_conn, syslog5424.LOG_DAEMON|syslog5424.LOG_WARNING, "custom-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	logger_err_conf := syslog.Channel(syslog5424.LOG_ERR).Logger("ERR : ")

	logger_err_conf.Print("doing some stuff")
	logger_err_conf.Print("doing anoter stuff")
	logger_err_conf.Print("doing a last stuff")

	time.Sleep(time.Second)

	sl_conn.End()
}



func server_custom(wg *sync.WaitGroup, mutex *sync.Mutex)  {
	wg.Add(1)
	defer wg.Done()

	collect, err	:= (syslog5424.Collector{
		QueueLen:	100,
	}).Collect("unix", TEST_SOCKET2, new(syslog5424.T_RFC5426) )
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

		fmt.Printf("%s\n", msg.String() )
	}

	time.Sleep(time.Second)

	collect.End()
}
