package extras4assetus

import "github.com/sneat-co/assetus/backend/const4assetus"

// DocTypeField is the per-field rule of a document-type validation schema,
// ported from the legacy frontend IDocTypeField.
type DocTypeField struct {
	Type     string // "str" | "int" | "date"
	Required bool
	Exclude  bool
	Max      int
	Min      int
}

// DocTypeStandardFields is the set of standard fields a document type may
// constrain, ported from the legacy frontend IDocTypeStandardFields.
type DocTypeStandardFields struct {
	Title     *DocTypeField
	Number    *DocTypeField
	IssuedBy  *DocTypeField
	IssuedOn  *DocTypeField
	ValidTill *DocTypeField
	Members   *DocTypeField
}

// DocTypeDef is a document-type definition with its validation schema, ported
// from the legacy frontend DocTypeDef.
type DocTypeDef struct {
	ID     const4assetus.Type
	Fields DocTypeStandardFields
}

// standardDocTypesByID is the per-doc-type validation schema ported from the
// legacy frontend standardDocTypesByID. Keys are AssetDocumentType IDs.
// Passport (and driving license) require number + validity; birth/marriage
// certificates require number + issuedOn and exclude validity.
var standardDocTypesByID = map[const4assetus.Type]DocTypeDef{
	"other": {
		ID:     "other",
		Fields: DocTypeStandardFields{Title: &DocTypeField{Required: true}},
	},
	const4assetus.TypeDocumentPassport: {
		ID: const4assetus.TypeDocumentPassport,
		Fields: DocTypeStandardFields{
			Number:    &DocTypeField{Required: true},
			ValidTill: &DocTypeField{Required: true},
			Members:   &DocTypeField{Max: 1},
		},
	},
	const4assetus.TypeDocumentDrivingLicense: {
		ID: const4assetus.TypeDocumentDrivingLicense,
		Fields: DocTypeStandardFields{
			Number:    &DocTypeField{Required: true},
			ValidTill: &DocTypeField{Required: true},
			Members:   &DocTypeField{Max: 1},
		},
	},
	const4assetus.TypeDocumentBirthCert: {
		ID: const4assetus.TypeDocumentBirthCert,
		Fields: DocTypeStandardFields{
			Number:    &DocTypeField{Required: true},
			IssuedBy:  &DocTypeField{},
			IssuedOn:  &DocTypeField{Required: true},
			ValidTill: &DocTypeField{Exclude: true},
			Members:   &DocTypeField{Max: 1},
		},
	},
	const4assetus.TypeDocumentMarriageCert: {
		ID: const4assetus.TypeDocumentMarriageCert,
		Fields: DocTypeStandardFields{
			Number:    &DocTypeField{Required: true},
			IssuedBy:  &DocTypeField{},
			IssuedOn:  &DocTypeField{Required: true},
			ValidTill: &DocTypeField{Exclude: true},
			Members:   &DocTypeField{Max: 2},
		},
	},
}

// DocTypeSchema returns the validation schema for a document type, or false if
// the document type has no standard schema.
func DocTypeSchema(docType const4assetus.Type) (DocTypeDef, bool) {
	def, ok := standardDocTypesByID[docType]
	return def, ok
}
