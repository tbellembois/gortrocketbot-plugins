package ldap

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/tbellembois/gortrocketbot/rocket"
	"golang.org/x/text/language"
	"gopkg.in/ldap.v2"

	// localization
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	err         error
	ldapc       *ldap.Conn
	ldapr       *ldap.SearchResult
	ldapmaxattr int

	// ldap serveur url - example: ldap.foo.com
	ldapServerURL = os.Getenv("ROCKETP_LDAP_SERVERURL")
	// ldap server port - example: 389
	ldapServerPort = os.Getenv("ROCKETP_LDAP_SERVERPORT")
	// ldap search base - example: dc=foo,dc=com
	ldapServerBase = os.Getenv("ROCKETP_LDAP_SERVERBASE")
	// ldap search filter - example: (&(cn=*%s*)(|(customAttr=0)(customAttr=9)))
	// "%s" will be replaced by the user search
	ldapSearchFilter = os.Getenv("ROCKETP_LDAP_SEARCHFILTER")
	// number of max results to display
	ldapMaxResults = os.Getenv("ROCKETP_LDAP_MAXRESULTS")
	// Rocket.Chat output format - example: %s :e-mail: %s :telephone_receiver: %s\n
	// "%s" will be replaced by the ldap attributes
	ldapResultFormat = os.Getenv("ROCKETP_LDAP_RESULTFORMAT")
	// ldap attributes to retrieve
	ldapAttributes = []string{"cn", "mail", "telephoneNumber"}

	// i18n
	localizer *i18n.Localizer
	bundle    *i18n.Bundle
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
		return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "error", PluralCount: 1})
	}

	if len(ldapr.Entries) == 0 {
		return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "noentry", PluralCount: 1})
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
		result += "\n"
	}

	return result
}

func init() {

	fmt.Println(fmt.Sprintf("ldapServerURL: %s", ldapServerURL))
	fmt.Println(fmt.Sprintf("ldapServerPort: %s", ldapServerPort))
	fmt.Println(fmt.Sprintf("ldapServerBase: %s", ldapServerBase))
	fmt.Println(fmt.Sprintf("ldapSearchFilter: %s", ldapSearchFilter))
	fmt.Println(fmt.Sprintf("ldapMaxResults: %s", ldapMaxResults))
	fmt.Println(fmt.Sprintf("ldapResultFormat: %s", ldapResultFormat))

	// converting ldap max attributes
	if ldapmaxattr, err = strconv.Atoi(ldapMaxResults); err != nil {
		log.Panic(err.Error())
	}

	// i18n initialization
	bundle = &i18n.Bundle{DefaultLanguage: language.Make(os.Getenv("ROCKETP_LDAP_LANGUAGE"))}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes(LOCALES_EN, "en.toml")
	bundle.MustParseMessageFileBytes(LOCALES_FR, "fr.toml")

	localizer = i18n.NewLocalizer(bundle)

	help := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "help", PluralCount: 1})
	rocket.RegisterPlugin(rocket.Plugin{
		Name:        "tel",
		CommandFunc: tel,
		Args:        []string{},
		Help:        help,
		IsAllowed: func(models.User) bool {
			return true
		},
	})
}
