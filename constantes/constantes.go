package constantes

import "time"

//Adresses et ports
const MulticastAddr = "224.0.0.1:6666"
const ListeningServerPort = ":6667"
const ListeningClientPort = ":6668"

//Attente
const AttenteK time.Duration = 3
const Min int = 4
const Max int = 60

//Types de message
const SYNC uint8 = 0
const FOLLOW_UP uint8 = 1
const DELAY_REQUEST uint8 = 2
const DELAY_RESPONSE uint8 = 3

//Paramétrage pour simulation
const DeriveHorloge int8 = 0
const DelaisTransmission int8 = 0
