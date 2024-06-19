package main

import "phaidra-connect/domain"

func hydratePhaidraObject(conf config, metadata domain.PhaidraMetadata) (*domain.PhaidraMetadataTemplate, error) {
	oefos, err := searchPhaidraOefos(conf, metadata.OefosId)
	if err != nil {
		return nil, err
	}

	orgUnit, err := searchPhaidraOrgUnit(conf, metadata.OrgUnitId)
	if err != nil {
		return nil, err
	}

	meta := new(domain.PhaidraMetadataTemplate)
	meta.Title = metadata.Title
	meta.Description = metadata.Description
	meta.ResourceLink = metadata.ResourceLink
	meta.Author = metadata.Author
	meta.Keywords = metadata.Keywords
	meta.Oefos = oefos
	meta.OrgUnit = orgUnit

	return meta, nil
}
