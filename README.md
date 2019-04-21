# gortrocketbot-plugins

Gortrocketbot plugin repository.

## hello

A sample plugin.

## ldap

Search users in an LDAP directory.

LDAP parameters are retrieved by environment variables:
```bash
export LDAP_SERVERURL="ldap.foo.com"
export LDAP_SERVERPORT="389"
export LDAP_SERVERBASE="dc=foo,dc=com"
export LDAP_SEARCHFILTER="(&(cn=*%s*)(|(customAttr=0)(customAttr=9)))"
export LDAP_MAXRESULTS="10"
export LDAP_RESULTFORMAT="%s :e-mail: %s :telephone_receiver: %s"
```

Displayed information are:
- `cn`
- `mail`
- `telephoneNumber`