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
	//On se connecte en écriture sur l'adresse de multicast
	conn, err := net.Dial("udp", constantes.MulticastAddr) // write on port
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var id uint8 = 0
	//Ecriture infinie
	for {
		//Envoi du SYNC
		sync := message.Message{
			Genre: constantes.SYNC,
			Id:    id,
		}
		tmaster := time.Now().UnixNano()
		message.SendMessage(sync.SimpleString(), conn)
		fmt.Printf("Send SYNC %s\n", time.Now())
		//Envoi du FOLLOW_UP
		follow_up := message.Message{
			Genre: constantes.FOLLOW_UP,
			Id:    id,
			Temps: tmaster,
		}
		message.SendMessage(follow_up.String(), conn)
		fmt.Printf("Send FOLLOW_UP %s\n", time.Now())
		id++
		time.Sleep(constantes.AttenteK * time.Second)
	}
}

//Attente de réception d'un DELAY_REQUEST
func selfListener() {
	//Ecoute sur soi
	conn, err := net.ListenPacket("udp", constantes.ListeningServerPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	//Ecoute infinie
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		//Pour chaque message
		for s.Scan() {
			mess := message.CreateMessage(s.Text())
			fmt.Printf("%s received from %v %s\n", s.Text(), addr, time.Now())
			switch mess.Genre {
				case constantes.DELAY_REQUEST:{
					fmt.Printf("Type DELAY_REQUEST\n")
					go delayResponseSender(mess.Id, addr.String())
				}
				default:{
					fmt.Printf("Unknown operation has been received.\n")
				}
			}
		}
	}
}

//On répond au esclave après un DELAY_REQUEST
func delayResponseSender(id uint8, addr string){
	//On crée la connexion en écriture sur l'esclave passé en paramètre
	conn, err := net.Dial("udp", strings.Split(addr, ":")[0] + constantes.ListeningClientPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//Envoi de DELAY_REQUEST
	delayRequest := message.Message{
		Genre: constantes.DELAY_RESPONSE,
		Id:    id,
		Temps:  time.Now().UnixNano(),
	}
	message.SendMessage(delayRequest.String(), conn)
	fmt.Printf("Send DELAY_RESPONSE %s\n", time.Now())
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}