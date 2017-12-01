package syslog5424

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const burstSocket string = "./test-burst.socket"
const burstMessage string = "doing some stuff"
const burstPacket string = "<27>1 2014-12-20T14:04:00Z localhost client-app 1234 - - ERR : doing some stuff"
const burstCount int = 100000

type burstTestOk struct {
	sock      string
	network   string
	transport Transport
}

//*
func Test_Burst(t *testing.T) {
	seq := []burstTestOk{
		{"u5425", "unix", T_RFC5425},
		{"ulf", "unix", T_LFENDED},
		{"uzero", "unix", T_ZEROENDED},
		{"dlf", "unixgram", T_LFENDED},
		{"dzero", "unixgram", T_ZEROENDED},
		{"d5425", "unixgram", T_RFC5425},
	}

	for _, s := range seq {
		t.Logf("burst on [%s] [%s]", s.network, s.transport.String())
		err := burst(burstSocket+s.sock, s.network, s.transport)
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

	err = nil
	Now = func() time.Time {
		t, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2014-12-20T14:04:00Z", time.UTC)
		return t
	}

	wg := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	mutex.Lock()
	wg.Add(2)
	go serverBurst(wg, mutex, sock, n, t, burstCount)
	go clientBurst(wg, mutex, sock, n, t, burstCount+100)

	wg.Wait()
	mutex.Unlock()

	return
}

func clientBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport, count int) {
	defer wg.Done()

	// waiting the creation of the socket
	mutex.Lock()
	slConn, chanErr, err := (Dialer{
		FlushDelay: 100 * time.Millisecond,
	}).Dial(n, sock, t)
	if err != nil {
		log.Fatalf("client Dial %q", err)
	}
	defer slConn.End()

	go func() {
		if err := <-chanErr; err != nil {
			log.Fatalf("client chanErr %q", err)
		}
	}()

	syslog, err := New(slConn, LOG_DAEMON|LOG_WARNING, "client-app")
	if err != nil {
		log.Fatalf("client New %q", err)
	}
	syslog.TestMode()

	loggerErrorConf := syslog.Channel(LOG_ERR).Logger("ERR : ")

	for i := 0; i < count; i++ {
		loggerErrorConf.Print(burstMessage)
	}

}

func serverBurst(wg *sync.WaitGroup, mutex *sync.Mutex, sock, n string, t Transport, count int) {
	defer wg.Done()

	listener, err := GuessListener(n, sock)
	if err != nil {
		log.Fatalf("server Collect %q", err)
	}

	collect, chanErr := NewReceiver(listener, 100, t)
	defer collect.End()

	go func() {
		if err := <-chanErr; err != nil {
			log.Fatalf("client chanErr %q", err)
		}
	}()

	// socket is created
	mutex.Unlock()

	for i := 0; i < count; i++ {
		msg, err, _ := collect.Receive()
		if err != nil {
			log.Fatalf("server receive %q", err)
		}
		if msg.String() != burstPacket {
			log.Fatalf("server got %q expected %q", msg, burstPacket)
		}
	}
}

func Benchmark_Burst(b *testing.B) {
	sock := burstSocket + "-bench"
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
