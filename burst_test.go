package syslog5424

import (
	"os"
	"log"
	"fmt"
	"time"
	"sync"
	"errors"
	"testing"
)


const	BURST_SOCKET	string	= "./test-burst.socket"
const	BURST_MESSAGE	string	= "doing some stuff"
const	BURST_PACKET	string	= "<27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing some stuff"
const	BURST_COUNT	int	= 1000

type	burst_tok struct {
	sock		string
	network		string
	transport	Transport
}


func Test_Burst(t *testing.T)  {
	seq := []burst_tok{
		{ "ulf",	"unix",		new(T_LFENDED)		},
		{ "uzero",	"unix",		new(T_ZEROENDED)	},
		{ "u5426",	"unix",		new(T_RFC5426)		},
		{ "dlf",	"unixgram",	new(T_LFENDED)		},
		{ "dzero",	"unixgram",	new(T_ZEROENDED)	},
		{ "d5426",	"unixgram",	new(T_RFC5426)		},
	}

	for _,s := range seq {
		//log.Printf("burst on [%s] [%s]\n", s.network, s.transport.String() )
		t.Logf("burst on [%s] [%s]", s.network, s.transport.String() )
		err := burst( BURST_SOCKET+s.sock, s.network, s.transport )
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}


func burst(sock,n string, t Transport) (err error) {
	defer os.Remove(sock)
	os.Remove(sock)

	wg	:= new(sync.WaitGroup)
	mutex	:= new(sync.Mutex)

	err	= nil
	mutex.Lock()

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}


	wg.Add(2)
	go serverBurst(wg, mutex, sock, n, t )
	go clientBurst(wg, mutex, sock, n, t )

	wg.Wait()
	mutex.Unlock()

	return
}


func clientBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport)  {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	sl_conn,err := (Dialer{
		QueueLen:	100,
		FlushDelay:	1*time.Second,
	}).Dial(n, sock, t )
	if err != nil {
		log.Fatal(err)
	}

	syslog, err := New(sl_conn, LOG_DAEMON|LOG_WARNING, "client-app")
	if err != nil {
		log.Fatal(err)
	}
	syslog.TestMode()

	logger_err_conf := syslog.Channel(LOG_ERR).Logger("ERR : ")

	for i:=0 ; i < BURST_COUNT+100 ; i++  {
		logger_err_conf.Print(BURST_MESSAGE)
	}

	sl_conn.End()
}



func serverBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport)  {
	defer wg.Done()

	collect, err	:= (Collector{
		QueueLen:	100,
	}).Collect(n, sock, t)
	if err != nil {
		log.Fatal(err)
	}

	// socket is created
	mutex.Unlock()

	for i:=0 ; i < BURST_COUNT ; i++  {
		msg, err, _ := collect.Receive()
		if err != nil {
			log.Fatal(err)
		}
		if msg.String() != BURST_PACKET {
			panic(errors.New(fmt.Sprintf("  got : [%s]", msg.String() )))
		}
	}

	collect.End()
}
