* ENGLISH

This script is very complex. Please read the instructions carefully.
Just change the variables between "SETTINGS START" and "SETTINGS END".

If you want to fly normal expeditions, you only have to set the following variables:

* home          Enter here your planet/moon, on which you want to carry out your expeditions, in "" 
                For example, enter "P:1:234:5" for the planet and "M:1:234:5" for the moon.

* radius        Starting from your starting planet, enter a radius in which the expeditions will be carried out. 
                With radius = 0 it remains in the same system. 

* rankOnePoints Enter the points of the best player in your universe.

* maxSlotsUse   !IMPORTANT! Here you tell the script how many SLOTS (not Expedition Slots) it can use.

* loop          "true" means that the script will continue to send until the user stops it by hand. 
                "false" means that when the radius is processed, the script ends.

* expoTime      Explained very simply: Holding time of the expedition.

* useWave       "true" means that all maximum slots are sent to one coordinate and only then the expedition changes. 
                "false" means that they are sent to each coordinate only once.

* useTelegram   if you want to have information by telegram, set this value to "true" - otherwise to "false".

* smallCargo    "false" means that large transporters are preferred in the expedition. "true" means small transporters.

* usePathfinder "true" means that the script always sends a pathfinder to an expedition. Not with "false".
                We recommend that this should always be "true". (more profit)

* useReaper     "true" means, the script always sends a Reaper instead of a destroyer to an expedition. 
                With "false" a destroyer is sent.

* HeavyFighter  "true" means that a maximum of 3000 heavy fighters will be sent to each expedition. If false, no.

* selfShips     "true" means you set the ships yourself. "false" means they are calculated automatically.

* TeleID        TELEGRAM_CHAT_ID is the ID that's in your settings. 
                However, you can specify another ID by using a possible second account. (CloudHoster only)

* mineDebris    If "true" the script checks if a debris field is present and removes it if necessary. With "false" not.
                IMPORTANT! (currently only for discoverers) This is only possible if "useWave" is also on "true".

* ship          With this variable you enter your preferred fleet, which should always be sent. Note that "selfShips" must be set to "true".
                Don't worry: Too much is not possible. The script checks this and reduces the values if necessary.


* SPECIAL SETTINGS

* hartDebris    "true" means, the script calculates a fleet composition for possible more debris field during fights.
                In plain language: It sends all found small or large transporters in a relation with light or heavy fighters. (see "LFandLC")
                "false" takes over the previous settings

* LFandLC       "true" means light fighters and small transporters are sent in a relation of 1:2
                "false" means heavy fighters and large transporters are sent in a relation of 1:1.5.

Everything else is done and checked by the script.
In this sense I wish you a lot of fun and success with this script!





* DEUTSCH

Dieses Skript ist sehr komplex. Bitte lies dir die Anleitung genau durch.
Ändere nur die Variablen zwischen "SETTINGS START" und "SETTINGS END"

Wenn du ganz normale Expeditionen fliegen willst, musst du nur folgende Variablen einstellen:

* home          Gib hier deinen Planeten/Mond, auf dem du deine Expeditionen durchführen willst, in "" ein. 
                Gib für den Planeten z.B. "P:1:234:5" ein und für den Mond "M:1:234:5".

* radius        Gib ausgehend von deinem Startplaneten einen Radius an, indem die Expeditionen durchgeführt werden. 
                Bei radius = 0 bleibt sie im selben System. 

* rankOnePoints Gib die Punkte des besten Spielers in deinem Universum an.

* maxSlotsUse   !WICHTIG! Hier sagts du dem Skript wie viele SLOTS (nicht Expeditions Slots) es maximal verwenden darf.

* loop          "true" bedeutet, dass das Skript immer weitersendet, bis der Benutzer es per Hand stoppt. 
                "false" bedeutet, wenn der Radius abgearbeitet wurde, endet das Skript.

* expoTime      Ganz einfach erklärt: Haltezeit der Expedition.

* useWave       "true" bedeutet, dass alle maximalen Slots zu einer Koordinate gesendet werden und erst dann gewechselt wird. 
                "false" bedeutet, dass sie zu jeder Koordinate nur einmal gesendet werden.

* useTelegram   wenn du Informationen per Telegram haben willst, setze diesen Wert auf "true" - andernfalls auf "false".

* smallCargo    "false" bedeutet, dass große Transporter in der Expedition bevorzugt wird. "true" bedeutet kleine Transporter.

* usePathfinder "true" bedeutet, das Skript sendet zu einer Expedition immer einen Pathfinder mit. Bei "false" nicht.
                Dies sollte unserer Emfpehlung nach immer "true" sein. (mehr Gewinn)

* useReaper     "true" bedeutet, das Skript sendet zu einer Expedition immer einen Reaper anstelle eines Zerstörers. 
                Bei "false" wird ein Zerstörer gesendet.

* HeavyFighter  "true" bedeutet, dass zu jeder Expedition maximal 3000 schwere Jäger gesendet werden. Bei "false" nicht.

* selfShips     "true" bedeutet, du stellst die Schiffe selber ein. "false" bedeutet, sie werden automatisch berechnet.

* TeleID        TELEGRAM_CHAT_ID ist die ID, die in deinen Einstellungen steht. 
                Du kannst jedoch durch einen möglichen zweiten Account eine weitere ID angeben. (Nur CloudHoster)

* mineDebris    Bei "true" überprüft das Skript ob ein Trümmerfeld vorhanden ist und baut dieses gegebenenfalls ab. Bei "false" nicht.
                !WICHTIG! (aktuell nur für Entdecker) Dies ist nur möglich, wenn "useWave" auch auf "true" ist.

* ship          Bei dieser Variable gibst du deine bevorzugte Flotte ein, die immer gesendet werden soll. Achte hierbei "selfShips" muss auf "true" sein.
                Keine Angst: Zu viel geht nicht. Das Skript überprüft dies und setzt gegebenenfalls die Werte runter.


* SPEZIALEINSTELLUNGEN

* hartDebris    "true" bedeutet, das Skript kalkuliert eine Flottenzusammenstellung für eventuell mehr Trümmerfeld bei Kämpfen.
                Im Klartext: Es sendet alle gefundenen kleinen oder großen Transporter in einem Verhältnis mit leichten oder schweren Jägern. (siehe "LFandLC")
                "false" übernimmt die vorherigen Einstellungen

* LFandLC       "true" bedeutet leichte Jäger und kleine Transporter werden in einem Verhältnis von 1:2 gesendet.
                "false" bedeutet schwere Jäger und große Transporter werden in einem Verhältnis von 1:1,5 gesendet.

Alles Weitere übernimmt und überprüft das Skript.
In diesem Sinne wünsche ich dir viel Spaß und viel Erfolg mit diesem Skript!


