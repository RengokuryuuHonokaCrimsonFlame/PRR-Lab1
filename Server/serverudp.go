package Server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/message"
	"github.com/RengokuryuuHonokaCrimsonFlame/PRR-Lab1/Client"

)

// debut, OMIT
const multicastAddr = "224.0.0.1:6666"

const SYNC = 0;
const FOLLOW_UP = 1;
const DELAY_REQUEST = 2;
const DELAY_RESPONSE = 3;

func main() {
	go clientReader()
	go syncSender()
	conn, err := net.Dial("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
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

	//A faire dane la partie esclave
	if err = p.JoinGroup(interf, addr); err != nil { // listen on ip multicast
		log.Fatal(err)
	}
	//
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

func syncSender(){
	conn, err := net.ListenPacket("udp", multicastAddr) // listen on port
	if err != nil {
		log.Fatal(err)
	}
	var id uint8
	defer conn.Close()
	for{
		mess := message.Message{
			SYNC,id, nil,
		}
		broadcast(mess)
	}
}

func broadcast(send chan message.Message) {
	fmt.Printf("COMM: Broadcasting message to: %s%s\n")
	broadcastAddress, err := net.ResolveUDPAddr("udp", multicastAddr)
	printError("ResolvingUDPAddr in Broadcast failed.", err)
	localAddress, err := net.ResolveUDPAddr("udp", GetLocalIP())
	connection, err := net.DialUDP("udp", localAddress, broadcastAddress)
	printError("DialUDP in Broadcast failed.", err)
	defer connection.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for {
		message := <-send
		err := encoder.Encode(message)
		printError("Encode error in broadcast: ", err)
		_, err = connection.Write(buffer.Bytes())
		if err != nil {
			_, err = lConnection.Write(buffer.Bytes())
			printError("Write in broadcast localhost failed", err)
		}
		buffer.Reset()
	}
}

func printError(erreur string, errorType error) {
	if errorType != nil {
		fmt.Println(erreur)
		fmt.Println(errorType.Error())
	}
}


// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

