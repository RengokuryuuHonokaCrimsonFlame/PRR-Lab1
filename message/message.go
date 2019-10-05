package message

import (
	"fmt"
	"time"
)

type Message struct{
	genre uint8
	id uint8
	temps time.Time
}

func (m Message) String() string{
	return fmt.Sprintf("<%v, %v , %v>", m.genre, m.id, m.temps)
}