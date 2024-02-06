package models

type OrgMember struct {
	UserID      string `bson:"userid" json:"-"` // to Disable UserID from Json encoding  use json:"-"
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	AccessLevel string `bson:"accesslevel" json:"accesslevel"`
}

// Organization represents the organization document in MongoDB
type Organization struct {
	OrgID       string      `bson:"_id"          json:"_id"`
	Name        string      `bson:"name"         json:"name"`
	Description string      `bson:"description"  json:"description"`
	Members     []OrgMember `bson:"members"      json:"members"`
}

type OrganizationBind struct {
	Name        string      `bson:"name"         json:"name"`
	Description string      `bson:"description"  json:"description"`
	Members     []OrgMember `bson:"members"      json:"members"`
}

type OrganizationOnly struct {
	OrgID       string `bson:"_id"          json:"_id"`
	Name        string `bson:"name"         json:"name"`
	Description string `bson:"description"  json:"description"`
}
