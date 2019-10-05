package Client

import (
	"io"
	"log"
	"net"
	"os"
)

// debut, OMIT
const srvAddr = "127.0.0.1:6000"
const multicastAddr = "224.0.0.1:6666"

func main() {
	conn, err := net.Dial("udp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		mustCopy(os.Stdout, conn)
	}()
	mustCopy(conn, os.Stdin) // CTRL-D pour sortir
}

func getLocalIP() (string) {
	var localIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	return localIP
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
// fin, OMIT
