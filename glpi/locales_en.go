package glpi

var LOCALES_EN = []byte(`
[help]
	one = "Search an item in GLPI"

[commandName]
	one = "name"
[commandNameHelp]
	one = "` + "`glpi` " + "`name`" + ` [exact_name]"
[commandPname]
	one = "pname"
[commandPnameHelp]
	one = "` + "`glpi` " + "`pname`" + ` [part_of_name]"
[commandOtherserial]
	one = "serial"
[commandOtherserialHelp]
	one = "` + "`glpi` " + "`serial`" + ` [inventory_number_with_spaces]"
[commandIPAdresses]
	one = "ip"
[commandIPAdressesHelp]
	one = "` + "`glpi` " + "`ip`" + ` [ip_address]"

[badCommand]
	one = "bad command"
`)
