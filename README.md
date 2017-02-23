# TTK4145

Sanntidsprogrammering NTNU -> aka verdens beste lab <3

YO, gjengen. om noen av dere skulle ende opp med å gjøre noe før jeg kommer/uten meg så legger jeg inn en liten update her <3
 
Har laget en FSMbrox som jeg tror er ganske sweet, i tillegg laget en mappe "driver2" som passer til denne. Det "nye" i den driver mapen er egenlig bare "driver" filen. Her har jeg oppdatert til hva jeg føler er en "ryddigere" måte. Nå er alle events en struct "Event" som kommer på eventchannelen. Denne kan ha en av to typer "BUTTONPRESSED" eller "NEWFLOOR" avhengig av hva slags event. Også skrevet om så nå kommer det event følgende ganger:

Button: Når en button blir presset, dvs at den må slippes og så trykkes inn på nytt for å sende. Sjekker ikke om den "lyser" eller tilsvarende bare om den er holdt inne siden sist vi "loopet" forbi den knappen

NewFloor: Når den kommer til en etg. etter å ha vært i "etg" "-1". 

Fjernet også set dir o.l fra Elev_int (?) da dette blir dekket i case INITIALIZE i fsmen -> om vi ikke bruker denne FSMen så må det sees på!

i FSM brukes en bool for å si fra om vi er annkommet ny etg. --> kan gjøres smidigere tror jeg

- BROX