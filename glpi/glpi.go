package glpi

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"net/http"
	"net/url"

	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/gorilla/schema"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gortrocketbot/rocket"
	"golang.org/x/text/language"
)

// GLPI field id to search
// and its associated help
type glpiSearch struct {
	// GLPI field id
	searchFieldID int
	// command help
	searchHelp string
}

// GLPI search returned data
type glpiData struct {
	Name      string   `json:"1"`
	Entity    string   `json:"80"`
	Inventoty string   `json:"6"`
	Serial    string   `json:"5"`
	Location  string   `json:"3"`
	Person    string   `json:"70"`
	Mac       []string `json:"21"`
	IP        []string `json:"126"`
	Type      string   `json:"itemtype"`
}

// GLPI search response
type glpiResponse struct {
	// GLPI session token after authentication
	SessionToken string     `json:"session_token"`
	Data         []glpiData `json:"data"`
}

// GLPI search request
type glpiRequest struct {
	Criteria []glpiRequestCriterion `json:"criteria"`
}

type glpiRequestCriterion struct {
	Field      int    `json:"field" schema:"field"`
	SearchType string `json:"searchtype" schema:"searchtype"`
	Value      string `json:"value" schema:"value"`
}

// GLPI data need an unmarshaller
// because the Mac and IP returned
// are string or slice depending of
// the number of fields returned
func (g *glpiData) UnmarshalJSON(data []byte) error {

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// setting "fixed" fields
	if v["1"] != nil {
		g.Name = v["1"].(string)
	}
	if v["80"] != nil {
		g.Entity = v["80"].(string)
	}
	if v["5"] != nil {
		g.Serial = v["5"].(string)
	}
	if v["6"] != nil {
		g.Inventoty = v["6"].(string)
	}
	if v["3"] != nil {
		g.Location = v["3"].(string)
	}
	if v["70"] != nil {
		g.Person = v["70"].(string)
	}
	g.Type = v["itemtype"].(string)

	// setting Mac and IP
	if v["21"] != nil {
		tmac := reflect.TypeOf(v["21"])

		switch tmac.Kind() {
		case reflect.Slice:
			for _, m := range v["21"].([]interface{}) {
				g.Mac = append(g.Mac, m.(string))
			}
		case reflect.String:
			g.Mac = []string{v["21"].(string)}
		}
	}

	if v["126"] != nil {
		tip := reflect.TypeOf(v["126"])
		switch tip.Kind() {
		case reflect.Slice:
			for _, m := range v["126"].([]interface{}) {
				g.IP = append(g.Mac, m.(string))
			}
		case reflect.String:
			g.IP = []string{v["126"].(string)}
		}
	}

	return nil
}

var (
	// GLPI application token
	glpiServerURL = os.Getenv("ROCKETP_GLPI_SERVERURL")
	// GLPI application token
	glpiAppToken = os.Getenv("ROCKETP_GLPI_APPTOKEN")
	// GLPI user
	glpiUser = os.Getenv("ROCKETP_GLPI_USER")
	// GLPI password
	glpiPassword = os.Getenv("ROCKETP_GLPI_PASSWORD")
	// search fields
	glpiSearchFields = map[string]glpiSearch{}
	// allowed users to access the plugin
	allowedUsers = strings.Split(os.Getenv("ROCKETP_GLPI_ALLOWEDUSERS"), ",")

	// i18n
	localizer *i18n.Localizer
	bundle    *i18n.Bundle

	// gorilla URL values encoder
	encoder *schema.Encoder

	// allowedUsers = []string{
	// 	"adleite",
	// 	"adrcanau",
	// 	"aleozwal",
	// 	"aligros",
	// 	"alkajfas",
	// 	"almirand",
	// 	"anmahul",
	// 	"antdurif",
	// 	"auvolle",
	// 	"caaube",
	// 	"cedchapa",
	// 	"clgendra",
	// 	"cyfayada",
	// 	"dabertra",
	// 	"dabiassi",
	// 	"dadelon",
	// 	"dagrimbi",
	// 	"delroche",
	// 	"didchaba",
	// 	"dovignal",
	// 	"dsibot",
	// 	"erboudon",
	// 	"erhren",
	// 	"ertourai",
	// 	"famonsei",
	// 	"flcogner",
	// 	"framrein",
	// 	"frrabet",
	// 	"gidubois",
	// 	"grcazorl",
	// 	"grrouchi",
	// 	"guilranc",
	// 	"gulevel",
	// 	"hedano",
	// 	"hedutill",
	// 	"hevidil",
	// 	"islauren",
	// 	"jechartr",
	// 	"jejougou",
	// 	"jetruffo",
	// 	"jfhitier",
	// 	"jlrenaud",
	// 	"jmcondat",
	// 	"jubonnau",
	// 	"juguegan",
	// 	"lalemoin",
	// 	"loverger",
	// 	"maaubert",
	// 	"mabacaul",
	// 	"mabenelh",
	// 	"mabouchu",
	// 	"mapoulon",
	// 	"marchana",
	// 	"meperrea",
	// 	"miortigu",
	// 	"moachira",
	// 	"naaurel",
	// 	"nadgoue",
	// 	"niboucha",
	// 	"nichabas",
	// 	"olbourgu",
	// 	"oltouret",
	// 	"paelie",
	// 	"paudevli",
	// 	"phjourde",
	// 	"pibarido",
	// 	"piraynau",
	// 	"ppchapon",
	// 	"rerobinm",
	// 	"rerossig",
	// 	"revillet",
	// 	"sarivet",
	// 	"secolin",
	// 	"seperbal",
	// 	"seplanch",
	// 	"stvalleg",
	// 	"syberaud",
	// 	"sylepeti",
	// 	"thbellem",
	// 	"thsmith",
	// 	"thvitry",
	// 	"vemuller",
	// 	"vithery",
	// 	"wiguyot",
	// 	"yoaraiso",
	// }
)

func help() string {
	h := ""
	h += localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandNameHelp", PluralCount: 1}) + "\n"
	h += localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPnameHelp", PluralCount: 1}) + "\n"
	h += localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandOtherserialHelp", PluralCount: 1}) + "\n"
	h += localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandIPAdressesHelp", PluralCount: 1}) + "\n"
	return h
}

func glpi(args ...string) string {

	var (
		req        *http.Request
		resp       *http.Response
		err        error
		body       []byte
		glpiresp   glpiResponse
		glpireq    glpiRequest
		glpireqGet string
	)

	if len(args) < 2 {
		return help()
	}
	cmd := args[0]
	search := strings.Join(args[1:], " ")

	// initializing GLPI session
	if req, err = http.NewRequest("GET", glpiServerURL+"/apirest.php/initSession", nil); err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Token", glpiAppToken)
	req.SetBasicAuth(glpiUser, glpiPassword)
	client := &http.Client{}
	if resp, err = client.Do(req); err != nil {
		return err.Error()
	}

	// reading response body
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return err.Error()
	}

	// getting session token
	if err = json.Unmarshal(body, &glpiresp); err != nil {
		return err.Error()
	}

	// building the request
	switch cmd {
	case localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPname", PluralCount: 1}):
		glpireq = glpiRequest{
			Criteria: []glpiRequestCriterion{
				glpiRequestCriterion{
					Field:      glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPname", PluralCount: 1})].searchFieldID,
					SearchType: "contains",
					Value:      search,
				},
			},
		}
	case localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandName", PluralCount: 1}):
		glpireq = glpiRequest{
			Criteria: []glpiRequestCriterion{
				glpiRequestCriterion{
					Field:      glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandName", PluralCount: 1})].searchFieldID,
					SearchType: "contains",
					Value:      "^" + search + "$",
				},
			},
		}
	case localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandOtherserial", PluralCount: 1}):
		glpireq = glpiRequest{
			Criteria: []glpiRequestCriterion{
				glpiRequestCriterion{
					Field:      glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandOtherserial", PluralCount: 1})].searchFieldID,
					SearchType: "contains",
					Value:      search,
				},
			},
		}
	case localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandIPAdresses", PluralCount: 1}):
		glpireq = glpiRequest{
			Criteria: []glpiRequestCriterion{
				glpiRequestCriterion{
					Field:      glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandIPAdresses", PluralCount: 1})].searchFieldID,
					SearchType: "contains",
					Value:      search,
				},
			},
		}
	default:
		return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "badCommand", PluralCount: 1})

	}

	// we support only simple searches
	// building the common GET parameters
	params := url.Values{}
	params.Add("criteria[0][field]", strconv.Itoa(glpireq.Criteria[0].Field))
	params.Add("criteria[0][searchtype]", glpireq.Criteria[0].SearchType)
	params.Add("criteria[0][value]", glpireq.Criteria[0].Value)
	params.Add("range", "0-500")
	params.Add("forcedisplay[0]", "80")
	params.Add("forcedisplay[1]", "1")
	params.Add("forcedisplay[2]", "32")
	params.Add("forcedisplay[3]", "11")
	params.Add("forcedisplay[4]", "6")
	params.Add("forcedisplay[5]", "5")
	params.Add("forcedisplay[6]", "3")
	params.Add("forcedisplay[7]", "70")
	params.Add("forcedisplay[8]", "21")
	params.Add("forcedisplay[9]", "126")
	glpireqGet = params.Encode()

	// performing request
	if req, err = http.NewRequest("GET", glpiServerURL+"/apirest.php/search/AllAssets?"+glpireqGet, nil); err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Token", glpiAppToken)
	req.Header.Set("Session-Token", glpiresp.SessionToken)
	client = &http.Client{}
	if resp, err = client.Do(req); err != nil {
		return err.Error()
	}

	// reading response body
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return err.Error()
	}

	// getting data
	if err = json.Unmarshal(body, &glpiresp); err != nil {
		return err.Error()
	}

	// building response
	ret := ""
	if len(glpiresp.Data) == 1 {
		ret = glpiresp.Data[0].Type + "\n"
		ret += "*nom*: " + glpiresp.Data[0].Name + "\n"
		ret += "*invent*: " + glpiresp.Data[0].Inventoty + "\n"
		ret += "*ser*: " + glpiresp.Data[0].Serial + "\n"

		for i, ip := range glpiresp.Data[0].IP {
			ret += "*ip" + strconv.Itoa(i) + "* `" + ip + "`\n"
		}
		for i, mac := range glpiresp.Data[0].Mac {
			ret += "*mac" + strconv.Itoa(i) + "* `" + mac + "`\n"
		}
		ret += glpiresp.Data[0].Entity + "\n"
		ret += glpiresp.Data[0].Location + "\n"
		ret += glpiresp.Data[0].Person + "\n"
	} else {
		for _, d := range glpiresp.Data {
			ret += "*nom*: " + d.Name + "\n"
		}
	}

	return ret
}

func init() {

	fmt.Println(allowedUsers)

	// gorilla schema encoder initialization
	encoder = schema.NewEncoder()

	// i18n initialization
	bundle = &i18n.Bundle{DefaultLanguage: language.Make(os.Getenv("ROCKETP_GLPI_LANGUAGE"))}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes(LOCALES_EN, "en.toml")
	bundle.MustParseMessageFileBytes(LOCALES_FR, "fr.toml")

	localizer = i18n.NewLocalizer(bundle)

	// commands initialization
	glpiSearchFields = make(map[string]glpiSearch)
	glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandName", PluralCount: 1})] = glpiSearch{
		searchFieldID: 1,
		searchHelp:    localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandNameHelp", PluralCount: 1}),
	}
	glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPname", PluralCount: 1})] = glpiSearch{
		searchFieldID: 1,
		searchHelp:    localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandPnameHelp", PluralCount: 1}),
	}
	glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandOtherserial", PluralCount: 1})] = glpiSearch{
		searchFieldID: 6,
		searchHelp:    localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandOtherserialHelp", PluralCount: 1}),
	}
	glpiSearchFields[localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandIPAdresses", PluralCount: 1})] = glpiSearch{
		searchFieldID: 126,
		searchHelp:    localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "commandIPAdressesHelp", PluralCount: 1}),
	}

	help := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "help", PluralCount: 1})
	rocket.RegisterPlugin(rocket.Plugin{
		Name:        "glpi",
		CommandFunc: glpi,
		Args:        []string{},
		Help:        help,
		IsAllowed: func(u models.User) bool {
			for _, v := range allowedUsers {
				if v == u.UserName {
					return true
				}
			}
			return false
		},
	})
}
