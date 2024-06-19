package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"phaidra-connect/domain"
	"text/template"
)

var conn *nats.Conn

var fns = template.FuncMap{
	"minus": func(a, b int) int {
		return a - b
	},
}

type config struct {
	natsHost        string
	museumHost      string
	phaidraHost     string
	phaidraUser     string
	phaidraPassword string
	templateFile    string
	template        *template.Template
}

func main() {
	conf := config{
		natsHost:        os.Getenv("NATS_HOST"),
		museumHost:      os.Getenv("MUSEUM_HOST"),
		phaidraHost:     os.Getenv("PHAIDRA_HOST"),
		phaidraUser:     os.Getenv("PHAIDRA_USERNAME"),
		phaidraPassword: os.Getenv("PHAIDRA_PASSWORD"),
		templateFile:    os.Getenv("TEMPLATE_FILE"),
	}

	tmpl, err := template.New(conf.templateFile).
		Funcs(fns).
		ParseFiles(conf.templateFile)
	if err != nil {
		panic(err)
	}

	conf.template = tmpl

	err = createPhaidraObject(conf, domain.PhaidraMetadata{
		Title:        "This was created by museum",
		Description:  "Test",
		ResourceLink: "https://sandbox.phaidra.org/objects/1",
		Author: domain.PhaidraMetadataAuthor{
			FirstName: "Gustav",
			LastName:  "Gans",
		},
		Keywords: [][]domain.PhaidraMetadataKeyword{
			{
				{Value: "Test", Lang: domain.PhaidraMetadataKeywordLangDE},
				{Value: "Test", Lang: domain.PhaidraMetadataKeywordLangEN},
			},
		},
		OefosId:   "504017",
		OrgUnitId: "A495",
	})
	if err != nil {
		panic(err)
	}
	return

	conn, err = nats.Connect(conf.natsHost)
	if err != nil {
		panic(err)
	}

	_, err = conn.Subscribe("museum.exhibit.created", func(m *nats.Msg) {
		var data map[string]interface{}
		err = json.Unmarshal(m.Data, &data)
		if err != nil {
			panic(err)
		}

		// create http get request to museum api
		/*apiPath := "/exhibit/" + data["exhibitId"]

		request, err := http.NewRequest(http.MethodGet, conf.museumHost+apiPath, nil)
		if err != nil {
			panic(err)
		}

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		fmt.Println(res.Body)*/

		fmt.Println(data)
	})
	if err != nil {
		panic(err)
	}

	select {}
}
