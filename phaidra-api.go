package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"phaidra-connect/domain"
	"strings"
)

func createPhaidraObject(conf config, metadata domain.PhaidraMetadata) error {
	apiPath := "/api/resource/create"

	conf.log.Debug("creating object in Phaidra", zap.String("apiPath", apiPath))

	hydratedMetadata, err := hydratePhaidraObject(conf, metadata)
	if err != nil {
		return err
	}

	metadataBuf := new(bytes.Buffer)
	err = conf.template.Execute(metadataBuf, hydratedMetadata)
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

	conf.log.Debug("created object in Phaidra", zap.String("response", resBody.String()), zap.String("apiPath", apiPath))

	return nil
}

func searchPhaidraOefos(conf config, oefos string) (*domain.PhaidraOefosMetadata, error) {
	apiPath := "/api/vocabulary?uri=oefos2012"

	conf.log.Debug("retrieving OEFOS from Phaidra", zap.String("oefos", oefos), zap.String("apiPath", apiPath))

	request, err := http.NewRequest(http.MethodGet, conf.phaidraHost+apiPath, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(conf.phaidraUser, conf.phaidraPassword)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	resBody := new(bytes.Buffer)
	_, err = resBody.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error searching OEFOS in Phaidra: %s", resBody.String())
	}

	root := new(domain.PhaidraVocabularyRoot)
	err = json.Unmarshal(resBody.Bytes(), &root)
	if err != nil {
		return nil, err
	}

	vocabulary, deLabel, enLabel := root.TreeSearch(oefos)
	if vocabulary == nil {
		return nil, fmt.Errorf("OEFOS %s not found in Phaidra", oefos)
	}

	meta := new(domain.PhaidraOefosMetadata)
	meta.FullLabel = []domain.PhaidraMetadataKeyword{
		{Value: deLabel, Lang: domain.PhaidraMetadataKeywordLangDE},
		{Value: enLabel, Lang: domain.PhaidraMetadataKeywordLangEN},
	}
	meta.ExactMatch = vocabulary.Id
	meta.Notation = oefos
	meta.PrefLabel = []domain.PhaidraMetadataKeyword{
		{Value: vocabulary.Labels[domain.PhaidraMetadataKeywordLangDE], Lang: domain.PhaidraMetadataKeywordLangDE},
		{Value: vocabulary.Labels[domain.PhaidraMetadataKeywordLangEN], Lang: domain.PhaidraMetadataKeywordLangEN},
	}

	conf.log.Debug("retrieved OEFOS from Phaidra", zap.String("label", enLabel), zap.String("exactMatch", meta.ExactMatch), zap.String("oefos", oefos), zap.String("apiPath", apiPath))

	return meta, nil
}

func searchPhaidraOrgUnit(conf config, orgUnitId string) (*domain.PhaidraOrgUnitMetadata, error) {
	apiPath := "/api/directory/org_get_units"

	conf.log.Debug("retrieving org unit from Phaidra", zap.String("orgUnitId", orgUnitId), zap.String("apiPath", apiPath))

	request, err := http.NewRequest(http.MethodGet, conf.phaidraHost+apiPath, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(conf.phaidraUser, conf.phaidraPassword)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	resBody := new(bytes.Buffer)
	_, err = resBody.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error searching org unit in Phaidra: %s", resBody.String())
	}

	root := new(domain.PhaidraOrgUnitRoot)
	err = json.Unmarshal(resBody.Bytes(), &root)
	if err != nil {
		return nil, err
	}

	orgUnit := root.TreeSearch(orgUnitId)
	if orgUnit == nil {
		return nil, fmt.Errorf("org unit %s not found in Phaidra", orgUnitId)
	}

	meta := new(domain.PhaidraOrgUnitMetadata)
	meta.FullLabel = []domain.PhaidraMetadataKeyword{
		{Value: orgUnit.PrefLabel[domain.PhaidraMetadataKeywordLangDE], Lang: domain.PhaidraMetadataKeywordLangDE},
		{Value: orgUnit.PrefLabel[domain.PhaidraMetadataKeywordLangEN], Lang: domain.PhaidraMetadataKeywordLangEN},
	}
	meta.ExactMatch = strings.ReplaceAll(orgUnit.Id, "/", "\\/")

	conf.log.Debug("retrieved org unit from Phaidra", zap.String("label", orgUnit.PrefLabel[domain.PhaidraMetadataKeywordLangEN]), zap.String("exactMatch", meta.ExactMatch), zap.String("orgUnitId", orgUnitId), zap.String("apiPath", apiPath))
	return meta, nil
}
