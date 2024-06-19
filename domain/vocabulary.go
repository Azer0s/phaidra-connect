package domain

import (
	"strings"
)

type PhaidraVocabularyRoot struct {
	Alerts     []string            `json:"alerts"`
	Vocabulary []PhaidraVocabulary `json:"vocabulary"`
}

type PhaidraVocabulary struct {
	Id       string                                `json:"@id"`
	Children []PhaidraVocabulary                   `json:"children"`
	Notation []string                              `json:"skos:notation"`
	Labels   map[PhaidraMetadataKeywordLang]string `json:"skos:prefLabel"`
}

func (v PhaidraVocabularyRoot) TreeSearch(id string) (vocab *PhaidraVocabulary, deLabel string, enLabel string) {
	for _, vocabulary := range v.Vocabulary {
		if result, deLabelRes, enLabelRes := vocabulary.TreeSearch(id); result != nil {
			return result, "ÖFOS 2012 -- " + deLabelRes, "ÖFOS 2012 -- " + enLabelRes
		}
	}

	return nil, "", ""
}

func (v PhaidraVocabulary) TreeSearch(id string) (vocab *PhaidraVocabulary, deLabel string, enLabel string) {
	oefosId := strings.Split(v.Id, ":")
	if len(oefosId) > 1 && oefosId[1] == id {
		return &v,
			v.Labels[PhaidraMetadataKeywordLangDE] + " (" + oefosId[1] + ")",
			v.Labels[PhaidraMetadataKeywordLangEN] + " (" + oefosId[1] + ")"
	}

	for _, child := range v.Children {
		if result, deLabelRes, enLabelRes := child.TreeSearch(id); result != nil {
			return result,
				v.Labels[PhaidraMetadataKeywordLangDE] + " (" + oefosId[1] + ") -- " + deLabelRes,
				v.Labels[PhaidraMetadataKeywordLangEN] + " (" + oefosId[1] + ") -- " + enLabelRes
		}
	}

	return nil, "", ""
}
