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
	"sync"
	"time"
	"golang.org/x/net/ipv4"
)

// debut, OMIT
var syncId uint8 = 0
var ecart int64 = 0
var addrServer = ""
var delayId uint8 = 0
var delay int64 = 0
var tes int64 = 0
var mutexSync sync.Mutex
var mutexEcart sync.Mutex
var mutexAddrServer sync.Mutex
var mutexDelayId sync.Mutex
var mutexDelay sync.Mutex
var mutexTes sync.Mutex

func main() {
	go udpReader()
	go delayResponseReceiver()
	conn, err := net.Dial("udp", constantes.MulticastAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(conn, os.Stdin)
}

// milieu, OMIT
//Écoute sur l'adresse multicast
func udpReader() {
	//Connexion
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
	
	//Ecoute infinie
	for {
		n, addr, err := conn.ReadFrom(buf) // n, _, addr, err := p.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		
		//Lit ce qui a été reçu
		for s.Scan() {
			mess := message.CreateMessage(s.Text())
			fmt.Printf( "%s received from %v %s\n", mess, addr, time.Now())
			mutexSync.Lock()
			if mess.Id < syncId { //Si l'on reçoit des anciens message
				mutexSync.Unlock()
				fmt.Printf( "Ancien message reçu.\n")
			} else {
				mutexSync.Unlock()
				switch mess.Genre { //Si l'on reçoit des message valide
				case constantes.SYNC:{
						fmt.Printf("Type SYNC\n")
						mutexSync.Lock()
						syncId = mess.Id
						mutexSync.Unlock()
						mutexAddrServer.Lock()
						
						// Lors du premier SYNC on lance la fonction delayRequestSender
						if addrServer != addr.String(){
							addrServer = addr.String()
							go delayRequestSender(addr.String())
						}
						mutexAddrServer.Unlock()
					}
					case constantes.FOLLOW_UP:{
						mutexEcart.Lock()
						ecart = time.Now().UnixNano() - mess.Temps - int64(constantes.DeriveHorloge); // Simulation de dérive
						fmt.Printf("Type FOLLOW_UP écart de %d nano secondes\n", ecart)
						mutexEcart.Unlock()
					}
					default:{
						fmt.Printf("Unknown operation has been received.")
					}
				}
			}
		}
	}
}

//Envoie une réponse au serveur
func delayRequestSender(addr string){
	//Boucle tant que l'adresse du serveur est la même (au cas où l'on changerait de serveur)
	for getAddrServer() == addr {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(constantes.Max - constantes.Min + 1) +  constantes.Min //Attente aléatoire
		time.Sleep(time.Duration(r) * time.Second)
		
		//Connexion au serveur
		conn, err := net.Dial("udp", strings.Split(addr, ":")[0]+constantes.ListeningServerPort)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		
		//Incrément de l'ID et envoi de la requête
		mutexDelayId.Lock()
		delayId += 1
		mess := message.Message{
			Genre: constantes.DELAY_REQUEST,
			Id:    delayId,
		}
		mutexDelayId.Unlock()
		time.Sleep(constantes.DelaisTransmission * time.Second) // Simulation de delais
		fmt.Printf("Send DELAY_REQUEST %s\n", time.Now())
		message.SendMessage(mess.SimpleString(), conn)
		mutexTes.Lock()
		tes = time.Now().UnixNano() - int64(constantes.DeriveHorloge); // Simulation de dérive
		mutexTes.Unlock()
	}
}

//Reçoit la réponse du serveur
func delayResponseReceiver(){
	//On écoute sur soi
	conn, err := net.ListenPacket("udp", constantes.ListeningClientPort) // listen on port
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
		
		//Pour chaque message reçu
		for s.Scan() {
			mess := message.CreateMessage(s.Text())
			fmt.Printf("%s received from %v %s\n", mess, addr, time.Now())
			switch mess.Genre{
				case constantes.DELAY_RESPONSE:{
					fmt.Printf("Type DELAY_RESPONSE\n")
					mutexDelayId.Lock()
					
					if delayId == mess.Id{ //On s'assure que l'id reçu corresponde à l'id attendu
						mutexDelayId.Unlock()
						mutexDelay.Lock()
						mutexTes.Lock()
						delay = (mess.Temps - tes) / 2
						mutexTes.Unlock()
						fmt.Printf("delais de %d nano secondes\n", delay)
						mutexDelay.Unlock()
					}else{
						mutexDelayId.Unlock()
						fmt.Printf("Ids don't match %d vs %d\n", mess.Id, delayId)
					}
				}
				default:{
					fmt.Printf("Unknown operation has been received.")
				}
			}
		}
	}
}

//Retourne l'adresse du serveur actuel
func getAddrServer() string{
	mutexAddrServer.Lock()
	result := addrServer
	mutexAddrServer.Unlock()
	return result
}

// fin, OMIT
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}