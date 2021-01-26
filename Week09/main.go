package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/James-Ren/Go-001/tree/main/Week09/internal/comet"
)

func main() {
	startTCPServer()
}

func startTCPServer() {
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	if err != nil {
		log.Printf("Resolve Addr Error:%v", err)
		return
	}
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Printf("Resolve Addr Error:%v", err)
		return
	}
	defer l.Close()
	log.Println("Start TCP Server")
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Printf("listener Accept Error:%v", err)
			return
		}
		if err = conn.SetReadBuffer(4096); err != nil {
			log.Printf("conn SetReadBuffer err:%v", err)
			return
		}
		if err = conn.SetWriteBuffer(4096); err != nil {
			log.Printf("conn SetWriteBuffer err:%v", err)
			return
		}
		go serveConn(conn)
	}
}

func serveConn(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	hb := time.Minute //heartbeat interval
	conn.SetReadDeadline(time.Now().Add(hb))
	signal := make(chan *comet.Proto)
	go dispatchConn(conn, signal)
	for {
		p := &comet.Proto{}

		err := p.ReadTCP(reader)
		if err != nil {
			log.Printf("Read failed:%v", err)
			break
		}
		if p.Op == comet.OpHeartbeat {
			conn.SetReadDeadline(time.Now().Add(hb))
			p.Op = comet.OpHeartbeatReply
			p.Body = nil
			signal <- p
		} else if p.Op == comet.OpSendMsg {
			p.Op = comet.OpSendMsgReply
			signal <- p
		}
	}
	conn.Close()
	signal <- comet.ProtoFinish //通知写goroutine退出
	log.Printf("conn Read goroutine closed")
}

func dispatchConn(conn *net.TCPConn, signal chan *comet.Proto) {
	var finish bool
	writer := bufio.NewWriter(conn)
	for {
		p := <-signal
		switch p {
		case comet.ProtoFinish:
			finish = true
			goto failed
		default:
			if err := p.WriteTCP(writer); err != nil {
				log.Printf("Write Proto failed:%v", err)
				goto failed
			}
		}

	}
failed:
	conn.Close()
	//防止读goroutine block signal
	for !finish {
		p := <-signal
		if p == comet.ProtoFinish {
			finish = true
		}
	}
	log.Printf("conn Write goroutine closed")
}
