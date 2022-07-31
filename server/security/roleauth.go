package security

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	c "github.com/thoniwutr/schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

type RoleAuthenticator interface {
	CheckPermission(req *http.Request, organisationID string, supportedRoles ...model.RoleType) error
}

type roleAuth struct {
	db db.RoleDB
}

func NewRoleAuth(db db.RoleDB) *roleAuth {
	return &roleAuth{db: db}
}

func (ra *roleAuth) CheckPermission(req *http.Request, organisationID string, supportedRoles ...model.RoleType) error {
	userInfo, err := getUserInfo(req)
	if err != nil {
		return err
	}

	role, err := ra.db.GetRole(organisationID, userInfo.ID)
	if err != nil {
		return err
	}

	for _, r := range role.RoleTypes {
		for _, supportedRole := range supportedRoles {
			if r == supportedRole {
				return nil
			}
		}
	}

	return c.ErrNoPermission
}

// authUserInfo is a struct that hold authenticated user info from Google Cloud Endpoint ESP
type authUserInfo struct {
	Issuer string `json:"issuer"`
	ID     string `json:"id"`
	Email  string `json:"email"`
}

// getUserInfo extract authenticated user information written to request header by Google Cloud Endpoint ESP
func getUserInfo(r *http.Request) (*authUserInfo, error) {
	encodedInfo := r.Header.Get("X-Endpoint-API-UserInfo")

	if encodedInfo == "" {
		return nil, c.ErrUnauthorized
	}

	data, err := base64.URLEncoding.DecodeString(encodedInfo)
	if err != nil {
		return nil, err
	}

	userInfo := &authUserInfo{}
	err = json.Unmarshal(data, userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
