package constantes

import "time"

const MulticastAddr = "224.0.0.1:6666"
const ListeningPort = ":6667"

const AttenteK time.Duration = 3
const Min int8 = 4
const Max int8 = 60

const SYNC uint8 = 0
const FOLLOW_UP uint8 = 1
const DELAY_REQUEST uint8 = 2
const DELAY_RESPONSE uint8 = 3
