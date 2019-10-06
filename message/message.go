package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Message struct{
	Genre uint8
	Id 	  uint8
	Temps time.Time
}

func (m Message) String() string{
	return fmt.Sprintf("%v %v %v", m.Genre, m.Id, m.Temps)
}

func (m Message) SimpleString() string{
	return fmt.Sprintf("%v %v", m.Genre, m.Id)
}

func CreateMessage(s string) *Message{
	decompose := strings.Split(s, " ")
	genre, _ := strconv.ParseUint(decompose[0], 10, 8)

	id, _ := strconv.ParseUint(decompose[1], 10, 8)

	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, decompose[2])
	if err != nil && len(decompose) > 2{
		fmt.Println(err)
	}
	mess := Message{
		Genre: uint8(genre),
		Id: uint8(id),
		Temps: t,
	}
	return &mess
}