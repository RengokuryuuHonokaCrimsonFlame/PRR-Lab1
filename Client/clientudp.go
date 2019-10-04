package Client

import (
	"io"
	"log"
	"net"
	"os"
)

// debut, OMIT
const srvAddr = "127.0.0.1:6000"

func main() {
	conn, err := net.Dial("udp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		mustCopyClientSide(os.Stdout, conn)
	}()
	mustCopyClientSide(conn, os.Stdin) // CTRL-D pour sortir
}

func mustCopyClientSide(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
// fin, OMIT
