package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

var conn *nats.Conn

type config struct {
	natsHost   string
	museumHost string
}

func main() {
	conf := config{
		natsHost:   os.Getenv("NATS_HOST"),
		museumHost: os.Getenv("MUSEUM_HOST"),
	}

	var err error
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
