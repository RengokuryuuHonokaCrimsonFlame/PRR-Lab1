package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/constantes"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
	"golang.org/x/net/ipv4"
)

// debut, OMIT

var syncId uint8 = 0
var ecart int64 = 0
var addrServer string = ""
var delayId uint8 = 0
var delay int64 = 0
var tes = make(chan int64)

func main() {
	go udpReader()
	go delayResponceReceiver()
	conn, err := net.Dial("udp", constantes.MulticastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
}

// milieu, OMIT
//On écoute sur l'adresse multicast
func udpReader() {
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
			mess := message.CreateMessage(s.Text())
			fmt.Printf( "%s received from %v\n", mess, addr)
			if mess.Id < syncId {
				fmt.Printf( "Ancien message reçu.\n")
			}else{
				switch mess.Genre {
				case constantes.SYNC:{
						fmt.Printf("SYNC\n")
						syncId = mess.Id
						if addrServer != addr.String(){
							addrServer = addr.String()
							go delayRequestSender(addr.String())
						}
					}
					case constantes.FOLLOW_UP:{
						ecart = time.Now().UnixNano() - mess.Temps
						fmt.Printf("FOLLOW_UP écart de %d nano secondes\n", ecart)
					}
					default:{
						fmt.Printf("Unknown operation has been received.")
					}
				}
			}
		}
	}
}

//On envoie une réponse au serveur
func delayRequestSender(addr string){
	for addrServer == addr {
		rand.Seed(time.Now().UnixNano())
		r := constantes.Min//rand.Intn(constantes.Max - constantes.Min + 1) +  constantes.Min
		time.Sleep(time.Duration(r) * time.Second)
		conn, err := net.Dial("udp", strings.Split(addr, ":")[0]+constantes.ListeningServerPort)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		mess := message.Message{
			Genre: constantes.DELAY_REQUEST,
			Id:    delayId,
		}
		fmt.Println("Send DELAY_REQUEST\n")
		message.SendMessage(mess.SimpleString(), conn)
		tes <- time.Now().UnixNano()
		delayId += 1
	}
}

func delayResponceReceiver(){
	conn, err := net.ListenPacket("udp", constantes.ListeningClientPort) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			mess := message.CreateMessage(s.Text())
			fmt.Printf("%s received from %v\n", mess, addr)
			switch mess.Genre{
				case constantes.DELAY_RESPONSE:{
					fmt.Printf("DELAY_RESPONSE\n")
					if delayId == mess.Id{
						delay = (mess.Temps - <-tes) / 2
						fmt.Printf("delais de %d nano secondes\n", delay)
					}else{
						fmt.Printf("Ids don't match")
					}
				}
			}
		}
	}
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
