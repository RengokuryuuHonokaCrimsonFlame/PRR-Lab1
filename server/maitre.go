package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/constantes"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
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

// milieu, OMIT

//Envoie les message SYNC et FOLLOW_UP sur l'adresse multicast (en boucle)
func multicastSender() {
	conn, err := net.Dial("udp", constantes.MulticastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var id uint8 = 0
	for {
		sync := message.Message{
			Genre: constantes.SYNC,
			Id:    id,
		}
		tmaster := time.Now().UnixNano()
		message.SendMessage(sync.SimpleString(), conn)
		//fmt.Printf("Send SYNC\n")
		follow_up := message.Message{
			Genre: constantes.FOLLOW_UP,
			Id:    id,
			Temps: tmaster,
		}
		message.SendMessage(follow_up.String(), conn)
		//fmt.Printf("Send FOLLOW_UP\n")
		id++
		time.Sleep(constantes.AttenteK * time.Second)
	}
}

func selfListener() {
	conn, err := net.ListenPacket("udp", constantes.ListeningPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			fmt.Printf("ICI")
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		fmt.Printf("%s received from %v\n", s.Text(), addr)
		for s.Scan() {
			mess := message.CreateMessage(s.Text())
			fmt.Printf("%s received from %v\n", s.Text(), addr)
			switch mess.Genre {
				case constantes.DELAY_REQUEST:{
					fmt.Printf("DELAY_REQUEST\n")
					go delayResponseSender(mess.Id, addr.String())
				}
				default:{
					fmt.Printf("Unknown operation has been received.\n")
				}
			}
		}
	}
}

func delayResponseSender(id uint8, addr string){
	delayRequest := message.Message{
		Genre: constantes.DELAY_RESPONSE,
		Id:    id,
		Temps:  time.Now().UnixNano(),
	}
	conn, err := net.Dial("udp", strings.Split(addr, ":")[0] + ":6668") // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	message.SendMessage(delayRequest.String(), conn)
	fmt.Printf("Send DELAY_RESPONSE\n")
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
