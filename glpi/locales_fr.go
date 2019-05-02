package glpi

var LOCALES_FR = []byte(`
[help]
	one = "Recherche un element dans GLPI"

[commandName]
	one = "nom"
[commandNameHelp]
	one = "` + "`glpi` " + "`nom`" + ` [nom_exact]"
[commandPname]
	one = "pnom"
[commandPnameHelp]
	one = "` + "`glpi` " + "`pnom`" + ` [partie_du_nom] (retourne 50 résultats max)"
[commandOtherserial]
	one = "invent"
[commandOtherserialHelp]
	one = "` + "`glpi` " + "`invent`" + ` [numéro_inventaire_avec_espaces]"
[commandIPAdresses]
	one = "ip"
[commandIPAdressesHelp]
	one = "` + "`glpi` " + "`ip`" + ` [adresse_ip]"

[badCommand]
	one = "mauvaise commande"
`)
