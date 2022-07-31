package model

// Role is a struct for holding user role data
type Role struct {
	UserID         string
	UserEmail      string
	OrganisationID string

	RoleTypes []RoleType `datastore:",noindex"`
}

type RoleType string

const (
	RoleTypeOwner  RoleType = "owner"
	RoleTypeEditor RoleType = "editor"
	RoleTypeViewer RoleType = "viewer"
)
