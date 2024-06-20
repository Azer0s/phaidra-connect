package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"phaidra-connect/domain"
)

func getExhibitById(conf config, exhibitId string) (*domain.MuseumExhibit, error) {
	apiPath := "/api/exhibits/" + exhibitId

	conf.log.Debug("fetching exhibit from museum api", zap.String("exhibitId", exhibitId), zap.String("apiPath", apiPath))

	request, err := http.NewRequest(http.MethodGet, conf.museumHost+apiPath, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	exhibit := new(domain.MuseumExhibit)
	err = json.NewDecoder(res.Body).Decode(exhibit)
	if err != nil {
		return nil, err
	}

	conf.log.Debug("fetched exhibit from museum api", zap.String("exhibitId", exhibitId), zap.String("apiPath", apiPath))

	return exhibit, nil
}
