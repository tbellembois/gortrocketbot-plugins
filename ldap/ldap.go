package ldap

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tbellembois/gortrocketbot/rocket"
	"gopkg.in/ldap.v2"
)

var (
	err         error
	ldapc       *ldap.Conn
	ldapr       *ldap.SearchResult
	ldapmaxattr int

	// ldap serveur url - example: ldap.foo.com
	ldapServerURL = os.Getenv("LDAP_SERVERURL")
	// ldap server port - example: 389
	ldapServerPort = os.Getenv("LDAP_SERVERPORT")
	// ldap search base - example: dc=foo,dc=com
	ldapServerBase = os.Getenv("LDAP_SERVERBASE")
	// ldap search filter - example: (&(cn=*%s*)(|(customAttr=0)(customAttr=9)))
	// "%s" will be replaced by the user search
	ldapSearchFilter = os.Getenv("LDAP_SEARCHFILTER")
	// number of max results to display
	ldapMaxResults = os.Getenv("LDAP_MAXRESULTS")
	// Rocket.Chat output format - example: %s :e-mail: %s :telephone_receiver: %s\n
	// "%s" will be replaced by the ldap attributes
	ldapResultFormat = os.Getenv("LDAP_RESULTFORMAT")
	// ldap attributes to retrieve
	ldapAttributes = []string{"cn", "mail", "telephoneNumber"}
)

func tel(search ...string) string {

	result := ""

	// init ldap connection
	ldapc, err = ldap.Dial("tcp", fmt.Sprintf("%s:%s", ldapServerURL, ldapServerPort))
	if err != nil {
		log.Fatal(err)
	}
	defer ldapc.Close()

	// building ldap filter for user
	f := fmt.Sprintf(ldapSearchFilter, search[0])

	// prepare ldap search
	searchRequest := ldap.NewSearchRequest(
		ldapServerBase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		f,
		ldapAttributes,
		nil,
	)

	// ldap request
	if ldapr, err = ldapc.Search(searchRequest); err != nil {
		fmt.Println(err.Error())
		return "an error occured"
	}

	if len(ldapr.Entries) == 0 {
		return "no entry found"
	}

	// building result
	for i, e := range ldapr.Entries {

		if i > ldapmaxattr {
			break
		}

		var att []interface{}
		for _, a := range ldapAttributes {
			att = append(att, e.GetAttributeValue(a))
		}
		result += fmt.Sprintf(ldapResultFormat, att...)
	}

	return result
}

func init() {

	// converting ldap max attributes
	if ldapmaxattr, err = strconv.Atoi(ldapMaxResults); err != nil {
		log.Panic(err.Error())
	}

	rocket.RegisterPlugin(rocket.Plugin{
		Name:        "tel",
		CommandFunc: tel,
		Args:        []string{},
		Help:        "Search users telephone number",
	})
}
