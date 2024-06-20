package domain

type MuseumPhaidraMetadata struct {
	Title           string                     `json:"phaidra-title"`
	Description     string                     `json:"phaidra-description"`
	Creator         string                     `json:"phaidra-creator"`
	Keywords        [][]PhaidraMetadataKeyword `json:"phaidra-keywords"`
	AuthorFirstName string                     `json:"phaidra-author-firstname"`
	AuthorLastName  string                     `json:"phaidra-author-lastname"`
	OefosId         string                     `json:"phaidra-oefos"`
	OrgUnitId       string                     `json:"phaidra-orgunit"`
}

type MuseumExhibit struct {
	Meta MuseumPhaidraMetadata `json:"meta"`
}
