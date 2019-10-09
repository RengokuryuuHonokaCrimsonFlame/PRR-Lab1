package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/constantes"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/ipv4"

	_ "github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	_ "github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/constantes"
)

// debut, OMIT

var syncId uint8
var delayId uint8
var ecart int64
var addrServer string = ""

func main() {
	go udpReader()
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
	if addrServer != addr.String(){
		addrServer = addr.String()
		go sendDelayRequest(addr.String())
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
			switch mess.Genre {
				case constantes.SYNC:{
					if mess.Id < syncId{
						fmt.Printf( "SYNC ancien à ignorer\n")
					}else{
						fmt.Printf( "SYNC\n")
						syncId = mess.Id
					}
				}
				case constantes.FOLLOW_UP:{
					ecart = time.Now().UnixNano() - mess.Temps
					fmt.Printf( "FOLLOW_UP écart de %d nano secondes\n", ecart)
				}
				default:{
					fmt.Printf("Unknown operation has been received.")
				}
			}
		}
	}
}

//On envoie une réponse au serveur
func sendDelayRequest(addr string){
	for addrServer == addr {
		time.Sleep(10 * time.Second)
		conn, err := net.Dial("udp", strings.Split(addr, ":")[0]+constantes.ListeningPort)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		mess := message.Message{
			Genre: constantes.DELAY_REQUEST,
			Id:    syncId,
		}
		fmt.Println("Send DELAY_REQUEST")
		message.SendMessage(mess.SimpleString(), conn)
	}
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
