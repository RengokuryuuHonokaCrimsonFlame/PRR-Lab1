# PRR-Lab1

## Lancement du programme
...

## Vérification des résultats
...

## constantes.go
Ce fichier contient toutes les constantes utilisées par le programme.
On y retrouve l'adresse de multicast, les ports d'écoute, le délai d'attente k et ses bornes, 
ainsi que les constantes représentant les divers types de messages.

## message.go
Ce fichier contient la classe Message et toutes ses fonctions relatives (création et envoi)
Un message se définit par son type, son id et le moment auquel il a été créé.

## esclave.go
Lorsqu'un esclave est démarré, il commence par lancer une écoute continue sur l'adresse ip
définie pour le multicast afin de pouvoir recevoir les éventuels messages d'un maître.
Selon le type de message reçu, on effectue l'action correspondante :
- SYNC : l'esclave enregistre l'id qu'il a reçu du maître, de sorte à ce que s'il reçoit une requête retardée, il l'ignore.
Au premier SYNC reçu, il lancera la routine s'occupant d'envoyer les messages de type DELAY_REQUEST.
- FOLLOW_UP : l'esclave calcule l'écart de temps entre son horloge et celle du maître puis l'enregistre.

Suite à cela, l'esclave se met en écoute continue dans l'attente d'une réponse du maître.
Si cette réponse est de type DELAY_RESPONSE, ...

Enfin, <mustCopy> ...

## maitre.go
Lorsqu'un maître est démarré, il commence par lancer un envoi continu de messages de type
SYNC et FOLLOW_UP sur l'adresse ip définie pour le multicast.

Il se met ensuite en écoute continue dans l'attente d'un message de type DELAY_REQUEST.
Une fois ce message reçu, il répond à l'esclave en ...

Enfin, <mustCopy> ...
