package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"text/template"
)

var conn *nats.Conn

type config struct {
	natsHost        string
	museumHost      string
	phaidraHost     string
	phaidraUser     string
	phaidraPassword string
	templateFile    string
	template        *template.Template
}

type PhaidraMetadata struct {
	Title        string
	Description  string
	ResourceLink string
}

func createPhaidraObject(conf config, metadata PhaidraMetadata) error {
	//curl -X POST -u user:pass "https://sandbox.phaidra.org/api/resource/create" -F "metadata=@resource_metadata.json"

	// post request to phaidra api (placeholder metadata for now)
	apiPath := "/api/resource/create"

	// add form file (metadata as json-ld)
	buf := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(buf)
	defer func(multipartWriter *multipart.Writer) {
		err := multipartWriter.Close()
		if err != nil {
			panic(err)
		}
	}(multipartWriter)

	metadataBuf := new(bytes.Buffer)
	err := conf.template.Execute(metadataBuf, metadata)
	if err != nil {
		return err
	}

	tmpMetadataFile := path.Join(os.TempDir(), uuid.New().String()+".json")
	defer func() {
		err := os.Remove(tmpMetadataFile)
		if err != nil {
			panic(err)
		}
	}()

	file, err := os.Create(tmpMetadataFile)
	if err != nil {
		return err
	}

	_, err = file.WriteString(metadataBuf.String())
	if err != nil {
		return err
	}

	part, err := multipartWriter.CreateFormFile("metadata", file.Name())
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, conf.phaidraHost+apiPath, buf)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
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

	fmt.Println(resBody.String())
	return nil
}

func main() {
	conf := config{
		natsHost:        os.Getenv("NATS_HOST"),
		museumHost:      os.Getenv("MUSEUM_HOST"),
		phaidraHost:     os.Getenv("PHAIDRA_HOST"),
		phaidraUser:     os.Getenv("PHAIDRA_USER"),
		phaidraPassword: os.Getenv("PHAIDRA_PASSWORD"),
		templateFile:    os.Getenv("TEMPLATE_FILE"),
	}

	tmpl, err := template.ParseFiles(conf.templateFile)
	if err != nil {
		panic(err)
	}

	conf.template = tmpl

	err = createPhaidraObject(conf, PhaidraMetadata{
		Title:        "Test",
		Description:  "Test",
		ResourceLink: "https://sandbox.phaidra.org/objects/1",
	})
	if err != nil {
		return
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
