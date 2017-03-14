# TTK4145

Sanntidsprogrammering NTNU -> aka verdens beste lab <3

Problemer/issues:
* noen ganger når en heis er på så tar den CAB order to ganger -> altså den kommer til etg. stopper, lukker opp dører starter timer -> SÅ lukker den dører-> åpner dører igjen og først DA skrur den av lyset i CAB button..
* UPDATE!: nå tar de ikke hverandre sine jobber!endret: Merge(&AllCabOrders[elevID].CabButtons[floor], remote, elevID, LivingPeers, i consCab, mulig dette må gjøres i hall også... for sliten å tenkte etter. Det som er issue nå er v. f.eks. at en trekker ut IP kabel på "local" så får den kræsj med "index put of range" ett sted!
* VELDIG tilgjengelig for spm. etc. enten mail(oystein.brox@gmail.com) eller andre steder

Usefull CMDs:
* KJØR: "go build main.go && ./main --id=123
* "who" -> finner ut om noen er SSHet seg på, set etter "ipen" på slutten
* "ps -aux|grep < id fra who >" e.g. pts/0
* "kill -9 <prosess id fra sshd> (ene linje fra over)" e.g. 7427
* "grep -nr <word you want to find>"
simulate packet loss:
* init(only first time after reboot) "sudo modprobe sch_netem"
* sudo tc qdisc add dev eth0 root netem loss 15%" - for adding 15% !Only first time!
* sudo tc qdisc change dev eth0 root netem loss 20%" - for cahnging to 20% 


Things that mayby should be a function...?
* "ServidedOrder i FSM" aka vi clearer flere ganger? er dobble for loop i FSM flere steder...

Off Topic:
* NB! sjekk om den nå tåler chans uten buffer! !
* funskjonsnavn
* comments o.l-

Final CodeComplete
Main
* rekkefølge ting blir tatt inn i 
ConfigFile
* fjerne noen vi ikke bruker
ConsensusHall
* fjerne comments
ConsensusCab
* fjerne comments
ElevatorStates
* fjerne comments
Merge
* fjerne comments
Driver mappen
* bytte navn?
FSM
* fjerne comments
* kan noe puttes i funksjoner?
* butt for loop -> fikses i funksjon?
HalLReqAss
* fjerne comments
* bytte navn
