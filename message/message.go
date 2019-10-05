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