package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/constantes"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	"golang.org/x/net/ipv4"
)

// debut, OMIT

func main() {
	go multicastSender()
	go selfListener()
	conn, err := net.Dial("udp", constantes.MulticastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
}

//Envoie les message SYNC et FOLLOW_UP sur l'adresse multicast (en boucle)
func multicastSender() {
	conn, err := net.Dial("udp", constantes.MulticastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var id uint8
	id = 0
	for {
		sync := message.Message{
			Genre: constantes.SYNC,
			Id:    id,
		}
		tmaster := time.Now().UnixNano()
		message.SendMessage(sync.SimpleString(), conn)
		follow_up := message.Message{
			Genre: constantes.FOLLOW_UP,
			Id:    id,
			Temps: tmaster,
		}
		message.SendMessage(follow_up.String(), conn)
		id++
		time.Sleep(constantes * time.Second)
	}
}

func selfListener() {
	conn, err := net.ListenPacket("udp", constantes.ListeningPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, err := net.ResolveUDPAddr("udp", constantes.MulticastAddr)
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
			mess := message.CreateMessage(s.Text())
			fmt.Printf("%s received from %v\n", s.Text(), addr)
			switch mess.Genre {
				case constantes.DELAY_REQUEST:{

				}
			}
		}
	}
}

// milieu, OMIT
func clientReader() {
	conn, err := net.ListenPacket("udp", constantes.MulticastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	p := ipv4.NewPacketConn(conn) // convert to ipv4 packetConn
	addr, err := net.ResolveUDPAddr("udp", constantes.MulticastAddr)
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

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
