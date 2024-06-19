package domain

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
	OefosId      string
	OrgUnitId    string
}

type PhaidraOefosMetadata struct {
	FullLabel  []PhaidraMetadataKeyword
	ExactMatch string
	Notation   string
	PrefLabel  []PhaidraMetadataKeyword
}

type PhaidraOrgUnitMetadata struct {
	FullLabel  []PhaidraMetadataKeyword
	ExactMatch string
}

type PhaidraMetadataTemplate struct {
	Title        string
	Description  string
	ResourceLink string
	Author       PhaidraMetadataAuthor
	Keywords     [][]PhaidraMetadataKeyword
	Oefos        *PhaidraOefosMetadata
	OrgUnit      *PhaidraOrgUnitMetadata
}
