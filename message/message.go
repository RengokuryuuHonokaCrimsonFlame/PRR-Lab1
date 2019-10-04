package message

import (
	"fmt"
	"time"
)

type Message struct{
	genre uint8
	idMessage uint8
	temps time.Time
}

func (m Message) String() string{
	return fmt.Sprintf("<%v, %v , %v>", m.genre, m.idMessage, m.temps)
}