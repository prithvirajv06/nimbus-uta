package models

import (
	"time"

	"github.com/prithvirajv06/nimbus-uta/go/core/internal/utils"
)

type User struct {
	NIMB_ID      string       `bson:"nimb_id" json:"nimb_id"`
	Fname        string       `bson:"fname" json:"fname"`
	Lname        string       `bson:"lname" json:"lname"`
	Email        string       `bson:"email" json:"email"`
	Password     string       `bson:"password" json:"-"`
	Organization Organization `bson:"organization" json:"organization"`
	Role         Role         `bson:"role" json:"role"`
	Audit        Audit        `bson:"audit" json:"audit"`
}

type Organization struct {
	NIMB_ID string `bson:"nimb_id" json:"nimb_id"`
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Audit   Audit  `bson:"audit" json:"audit"`
}

type Role struct {
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

func NewOrganization(name, address string) *Organization {
	return &Organization{
		NIMB_ID: utils.GenerateNIMBID("N_ORG"),
		Name:    name,
		Address: address,
		Audit: Audit{Active: "ACTIVE",
			IsArchived: false,
			Version:    1, MinorVersion: 0,
			IsProdCandidate: false,
			CreatedAt:       time.Now().Unix(),
			CreatedBy:       "APP_REQUEST",
			ModifiedAt:      time.Now().Unix(),
			ModifiedBy:      "NONE"},
	}
}

func NewUser(fname, lname, email, password, organization string, role Role) *User {

	return &User{
		NIMB_ID:      utils.GenerateNIMBID("N_USER"),
		Fname:        fname,
		Lname:        lname,
		Email:        email,
		Password:     password,
		Organization: Organization{Name: organization},
		Role:         role,
		Audit: Audit{Active: "ACTIVE",
			IsArchived: false,
			Version:    1, MinorVersion: 0,
			IsProdCandidate: false,
			CreatedAt:       time.Now().Unix(),
			CreatedBy:       "APP_REQUEST",
			ModifiedAt:      time.Now().Unix(),
			ModifiedBy:      "NONE"},
	}
}

func (u *User) FullName() string {
	return u.Fname + " " + u.Lname
}

func (u *User) IsAdmin() bool {
	return u.Role.Name == "admin"
}

func (u *User) HasPermission(permission string) bool {
	for _, p := range u.Role.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func CreateDefaultAdminRole() Role {
	return Role{
		Name: "admin",
		Permissions: []string{
			"create_user",
			"delete_user",
			"update_user",
			"view_user",
			"manage_roles",
		},
	}
}
func CreateDefaultManagerRole() Role {
	return Role{
		Name: "manager",
		Permissions: []string{
			"create_user",
			"update_user",
			"view_user",
		},
	}
}

func CreateDefaultEditorRole() Role {
	return Role{
		Name: "editor",
		Permissions: []string{
			"update_user",
			"view_user",
		},
	}
}

func CreateDefaultUserRole() Role {
	return Role{
		Name: "user",
		Permissions: []string{
			"view_user",
		},
	}
}
