package syslog5424

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const BURST_SOCKET string = "./test-burst.socket"
const BURST_MESSAGE string = "doing some stuff"
const BURST_PACKET string = "<27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing some stuff"
const BURST_COUNT int = 100

type burst_tok struct {
	sock      string
	network   string
	transport Transport
}

//*
func Test_Burst(t *testing.T) {
	seq := []burst_tok{
		{"ulf", "unix", T_LFENDED},
		{"uzero", "unix", T_ZEROENDED},
		{"u5425", "unix", T_RFC5425},
		{"dlf", "unixgram", T_LFENDED},
		{"dzero", "unixgram", T_ZEROENDED},
		{"d5425", "unixgram", T_RFC5425},
	}

	for _, s := range seq {
		t.Logf("burst on [%s] [%s]", s.network, s.transport.String())
		err := burst(BURST_SOCKET+s.sock, s.network, s.transport)
		if err != nil {
			t.Logf("[%s] [%s] %v", s.network, s.transport.String(), err)
			t.Fail()
		}
	}
}

//*/
func burst(sock, n string, t Transport) (err error) {
	defer os.Remove(sock)
	os.Remove(sock)

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	err = nil
	mutex.Lock()

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	wg.Add(2)
	go serverBurst(wg, mutex, sock, n, t, BURST_COUNT)
	go clientBurst(wg, mutex, sock, n, t, BURST_COUNT+100)

	wg.Wait()
	mutex.Unlock()

	return
}

func clientBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport, count int) {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	sl_conn, chan_err, err := (Dialer{
		QueueLen:   100,
		FlushDelay: 50 * time.Millisecond,
	}).Dial(n, sock, t)
	if err != nil {
		log.Fatalf("client Dial %q", err)
	}
	defer sl_conn.End()

	go func() {
		if err := <-chan_err; err != nil {
			log.Fatalf("client chan_err %q", err)
		}
	}()

	syslog, err := New(sl_conn, LOG_DAEMON|LOG_WARNING, "client-app")
	if err != nil {
		log.Fatalf("client New %q", err)
	}
	syslog.TestMode()

	logger_err_conf := syslog.Channel(LOG_ERR).Logger("ERR : ")

	for i := 0; i < count; i++ {
		logger_err_conf.Print(BURST_MESSAGE)
	}

}

func serverBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport, count int) {
	defer wg.Done()

	collect, err := (Collector{
		QueueLen: 100,
	}).Collect(n, sock, t)
	if err != nil {
		log.Fatalf("server Collect %q", err)
	}
	defer collect.End()

	// socket is created
	mutex.Unlock()

	for i := 0; i < count; i++ {
		msg, err, _ := collect.Receive()
		if err != nil {
			log.Fatalf("server receive %q", err)
		}
		if msg.String() != BURST_PACKET {
			log.Fatalf("server got %q expected %q", msg, BURST_PACKET)
		}
	}

}

func Benchmark_Burst(b *testing.B) {
	sock := BURST_SOCKET + "-bench"
	defer os.Remove(sock)
	os.Remove(sock)

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	mutex.Lock()
	wg.Add(2)
	go clientBurst(wg, mutex, sock, "unix", T_RFC5425, b.N+100)

	b.ResetTimer()
	serverBurst(wg, mutex, sock, "unix", T_RFC5425, b.N)

	wg.Wait()
	mutex.Unlock()

	return

}
