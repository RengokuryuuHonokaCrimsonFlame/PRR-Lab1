package message

import (
	"fmt"
	"time"
)

type Message struct{
	Genre uint8
	Id 	  uint8
	Temps time.Time
}

func (m Message) String() string{
	return fmt.Sprintf("<%v, %v , %v>", m.Genre, m.Id, m.Temps)
}

func (m Message) SimpleString() string{
	return fmt.Sprintf("<%v, %v>", m.Genre, m.Id)
}

var MulticastAddr = "224.0.0.1:6666"
var ServerAddr = "127.0.0.1:6666"

var SYNC uint8 = 0
var FOLLOW_UP uint8 = 1
var DELAY_REQUEST uint8 = 2
var DELAY_RESPONSE uint8 = 3