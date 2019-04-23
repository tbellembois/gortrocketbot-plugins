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
