# PRR-Lab1
Crüll Loris, Rod Julien

## Lancement du programme
Tout d'abord, exécuter le fichier maitre.go pour démarrer un maître, puis exécuter le fichier esclave.go pour démarrer un esclave.

## constantes.go
Ce fichier contient toutes les constantes utilisées par le programme.
On y retrouve l'adresse de multicast, les ports d'écoute, le délai d'attente k et ses bornes, 
ainsi que les constantes représentant les divers types de messages.

|     Message    | Identifiant |
|:--------------:|:-----------:|
| SYNC           |      0      |
| FOLLOW_UP      |      1      |
| DELAY_REQUEST  |      2      |
| DELAY_RESPONSE |      3      |

## message.go
Ce fichier contient la classe Message et toutes ses fonctions relatives (création et envoi)
Un message se définit par son type (uint8), son id (uint8) et le moment auquel il a été créé (int64).
Lors de sa création, on incrémente son id correspondant avant de le lui attribuer et si besoin est, on lui ajoute la date courante.
Lors de l'envoi, le type du message est converti en string avant d'être re-converti en type Message à sa réception.

## esclave.go
Lorsqu'un esclave est démarré, il commence par lancer une écoute continue sur l'adresse ip
définie pour le multicast afin de pouvoir recevoir les éventuels messages d'un maître.
Selon le type de message reçu, on effectue l'action correspondante :
- SYNC : l'esclave enregistre l'id qu'il a reçu du maître, de sorte à ce que s'il reçoit une requête retardée, il l'ignore.
Au premier SYNC reçu, il lancera la routine s'occupant d'envoyer les messages de type DELAY_REQUEST.
- FOLLOW_UP : l'esclave calcule l'écart de temps entre son horloge et celle du maître puis l'enregistre.

Suite à cela, l'esclave se met en écoute continue dans l'attente d'une réponse du maître.
Si cette réponse est de type DELAY_RESPONSE, il enregistre le délai correspondant à la formule indiquée dans la donnée si l'id reçu correspond au dernier id envoyé.

## maitre.go
Lorsqu'un maître est démarré, il commence par lancer un envoi continu de messages de type
SYNC et FOLLOW_UP sur l'adresse ip définie pour le multicast.

Il se met en même temps en écoute continue sur lui-même dans l'attente d'un message de type DELAY_REQUEST.
Une fois ce message reçu, il répond à l'esclave avec un message de type DELAY_RESPONSE contenant l'id reçu et son temps actuel.

## Interprétation des résultats
En comparant les logs du maître et des différents esclaves, on s'aperçoit que les messages envoyés sont reçus des deux côtés.
Cela prouve le bon fonctionnement de notre implémentation, pour autant qu'aucun problème n'intervienne.
Si l'on redémarre le serveur, on constate dans les logs des clients que les messages SYNC et FOLLOW_UP sont ignorés.
Cela est dû au fait que les ids envoyés par le serveur repartent de 0 et sont considérés comme étant des anciens messages perdus.