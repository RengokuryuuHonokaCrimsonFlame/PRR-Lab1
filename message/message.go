package message

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Message struct{
	Genre uint8
	Id 	  uint8
	Temps int64
}

//Envoie un message
func (m Message) String() string{
	return fmt.Sprintf("%v %v %v", m.Genre, m.Id, m.Temps)
}

//Envoie un message sans temps
func (m Message) SimpleString() string{
	return fmt.Sprintf("%v %v", m.Genre, m.Id)
}

//Recrée le message a partir d'une string
func CreateMessage(s string) *Message{
	decompose := strings.Split(s, " ")
	genre, _ := strconv.ParseUint(decompose[0], 10, 8)

	id, _ := strconv.ParseUint(decompose[1], 10, 8)

	t := int64(time.Now().UnixNano()) / int64(time.Millisecond)
	if len(decompose) > 2 {
		t, _ = strconv.ParseInt(decompose[1], 10, 64)
	}

	mess := Message{
		Genre: uint8(genre),
		Id: uint8(id),
		Temps: t,
	}
	return &mess
}

//Pour envoyer un message
func SendMessage(message string, conn io.Writer){
	conn.Write([]byte(message))
}