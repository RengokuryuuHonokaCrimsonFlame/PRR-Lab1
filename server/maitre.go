package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

// debut, OMIT
const multicastAddr = "224.0.0.1:6666"

const SYNC = 0
const FOLLOW_UP = 1
const DELAY_REQUEST = 2
const DELAY_RESPONSE = 3

func main() {
	go multicastSender()
	conn, err := net.Dial("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
}

func multicastSender() {
	conn, err := net.Dial("udp", multicastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var id uint8
	id = 0
	for{
		mess := message.Message{
			Genre:	SYNC,
			Id:		id,
		}
		sendMessage(mess, conn)
		id++
		time.Sleep(10 * time.Second)
	}
}

// milieu, OMIT
func clientReader() {
	conn, err := net.ListenPacket("udp", multicastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	var interf *net.Interface
	if runtime.GOOS == "darwin" {
		interf, _ = net.InterfaceByName("en0")
	}

	if err = p.JoinGroup(interf, addr); err != nil { // listen on ip multicast
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			fmt.Printf("%s from %v\n", s.Text(), addr)
		}
	}
}

func sendMessage(mess message.Message, conn io.Writer){
	conn.Write([]byte(mess.String()))
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
