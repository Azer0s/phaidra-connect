package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"net/http"
	"net/url"
	"os"
	"strings"
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

type PhaidraMetadataAuthor struct {
	FirstName string
	LastName  string
}

type PhaidraMetadataKeywordLang string

const (
	PhaidraMetadataKeywordLangDE PhaidraMetadataKeywordLang = "deu"
	PhaidraMetadataKeywordLangEN PhaidraMetadataKeywordLang = "eng"
)

type PhaidraMetadataKeyword struct {
	Value string
	Lang  PhaidraMetadataKeywordLang
}

type PhaidraMetadata struct {
	Title        string
	Description  string
	ResourceLink string
	Author       PhaidraMetadataAuthor
	Keywords     [][]PhaidraMetadataKeyword
}

func createPhaidraObject(conf config, metadata PhaidraMetadata) error {
	//curl -X POST -u user:pass "https://sandbox.phaidra.org/api/resource/create" -F "metadata=@resource_metadata.json"

	// post request to phaidra api (placeholder metadata for now)
	apiPath := "/api/resource/create"

	metadataBuf := new(bytes.Buffer)
	err := conf.template.Execute(metadataBuf, metadata)
	if err != nil {
		return err
	}

	formData := url.Values{
		"metadata": {metadataBuf.String()},
	}

	request, err := http.NewRequest(http.MethodPost, conf.phaidraHost+apiPath, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(conf.phaidraUser, conf.phaidraPassword)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	resBody := new(bytes.Buffer)
	_, err = resBody.ReadFrom(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating object in Phaidra: %s", resBody.String())
	}
	return nil
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

	err = createPhaidraObject(conf, PhaidraMetadata{
		Title:        "Test",
		Description:  "Test",
		ResourceLink: "https://sandbox.phaidra.org/objects/1",
		Author: PhaidraMetadataAuthor{
			FirstName: "Gustav",
			LastName:  "Gans",
		},
		Keywords: [][]PhaidraMetadataKeyword{
			{
				{Value: "Test", Lang: PhaidraMetadataKeywordLangDE},
				{Value: "Test", Lang: PhaidraMetadataKeywordLangEN},
			},
		},
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
