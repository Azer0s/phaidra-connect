package main

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
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
	log             *zap.Logger
}

func main() {
	var log *zap.Logger
	var err error
	if os.Getenv("RUNTIME_ENV") == "production" {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}
	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
			panic(err)
		}
	}(log)

	conf := config{
		natsHost:        os.Getenv("NATS_HOST"),
		museumHost:      os.Getenv("MUSEUM_HOST"),
		phaidraHost:     os.Getenv("PHAIDRA_HOST"),
		phaidraUser:     os.Getenv("PHAIDRA_USERNAME"),
		phaidraPassword: os.Getenv("PHAIDRA_PASSWORD"),
		templateFile:    os.Getenv("TEMPLATE_FILE"),
		log:             log,
	}

	tmpl, err := template.New(conf.templateFile).
		Funcs(fns).
		ParseFiles(conf.templateFile)
	if err != nil {
		log.Fatal("error parsing template file", zap.Error(err))
	}

	conf.template = tmpl

	conn, err = nats.Connect(conf.natsHost)
	if err != nil {
		panic(err)
	}

	_, err = conn.Subscribe("museum.exhibit.created", func(m *nats.Msg) {
		log.Info("received museum.exhibit.created event")

		cloudEvent := new(cloudevents.Event)
		err = json.Unmarshal(m.Data, &cloudEvent)
		if err != nil {
			log.Error("error unmarshalling cloud event", zap.Error(err))
			return
		}

		data := new(map[string]interface{})
		err = json.Unmarshal(cloudEvent.Data(), data)
		if err != nil {
			log.Error("error unmarshalling cloud event data", zap.Error(err))
			return
		}

		log.Debug("creating phaidra object", zap.String("exhibitId", (*data)["exhibitId"].(string)))

		exhibit, err := getExhibitById(conf, (*data)["exhibitId"].(string))
		if err != nil {
			log.Error("error getting exhibit by id", zap.Error(err))
			return
		}

		meta := domain.PhaidraMetadata{
			Title:        exhibit.Meta.Title,
			Description:  exhibit.Meta.Description,
			ResourceLink: conf.museumHost + "/exhibit/" + (*data)["exhibitId"].(string),
			Author: domain.PhaidraMetadataAuthor{
				FirstName: exhibit.Meta.AuthorFirstName,
				LastName:  exhibit.Meta.AuthorLastName,
			},
			Keywords:  exhibit.Meta.Keywords,
			OefosId:   exhibit.Meta.OefosId,
			OrgUnitId: exhibit.Meta.OrgUnitId,
		}

		err = createPhaidraObject(conf, meta)
		if err != nil {
			log.Error("error creating phaidra object", zap.Error(err))
			return
		}

		log.Info("created phaidra object")
	})
	if err != nil {
		panic(err)
	}

	select {}
}
