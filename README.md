# gortrocketbot-plugins

[Gortrocketbot](https://github.com/tbellembois/gortrocketbot) plugin repository.

## hello

A sample plugin.

Language is set by environment variable:
```bash
export ROCKETP_HELLO_LANGUAGE="fr"
```

## ldap

Search users in an LDAP directory.

Language and LDAP parameters are set by environment variables:
```bash
export ROCKETP_LDAP_LANGUAGE="fr"
export ROCKETP_LDAP_SERVERURL="ldap.foo.com"
export ROCKETP_LDAP_SERVERPORT="389"
export ROCKETP_LDAP_SERVERBASE="dc=foo,dc=com"
export ROCKETP_LDAP_SEARCHFILTER="(&(cn=*%s*)(|(customAttr=0)(customAttr=9)))"
export ROCKETP_LDAP_MAXRESULTS="10"
export ROCKETP_LDAP_RESULTFORMAT="%s :e-mail: %s :telephone_receiver: %s"
```

Displayed information are:
- `cn`
- `mail`
- `telephoneNumber`

## glpi

Search informations in GLPI.
This is not a ready to use plugin. It must be modified to everyone's needs.

Language and GLPI parameters are set by environment variables:
```bash
export ROCKETP_GLPI_LANGUAGE="fr"
export ROCKETP_GLPI_SERVERURL="https://glpi.foo.com/glpi"
export ROCKETP_GLPI_APPTOKEN="glpigeneratedapptoken"
export ROCKETP_GLPI_USER="rocket"
export ROCKETP_GLPI_PASSWORD="mypassword"
export ROCKETP_GLPI_ALLOWEDUSERS="john,jean,eliot"
```