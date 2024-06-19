package domain

type PhaidraOrgUnitRoot struct {
	Alerts   []string         `json:"alerts"`
	Status   int              `json:"status"`
	OrgUnits []PhaidraOrgUnit `json:"units"`
}

func (r PhaidraOrgUnitRoot) TreeSearch(id string) *PhaidraOrgUnit {
	for _, orgUnit := range r.OrgUnits {
		if result := orgUnit.TreeSearch(id); result != nil {
			return result
		}
	}

	return nil
}

type PhaidraOrgUnit struct {
	Id        string                                `json:"@id"`
	Type      string                                `json:"@type"`
	Notation  string                                `json:"skos:notation"`
	PrefLabel map[PhaidraMetadataKeywordLang]string `json:"skos:prefLabel"`
	SubUnits  []PhaidraOrgUnit                      `json:"subunits"`
}

func (o PhaidraOrgUnit) TreeSearch(id string) *PhaidraOrgUnit {
	if o.Notation == id {
		return &o
	}

	for _, subUnit := range o.SubUnits {
		if result := subUnit.TreeSearch(id); result != nil {
			return result
		}
	}

	return nil
}
