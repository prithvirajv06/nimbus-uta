package models

import (
	"github.com/gin-gonic/gin"
	"github.com/prithvirajv06/nimbus-uta/go/core/internal/utils"
)

type User struct {
	NIMB_ID      string       `bson:"nimb_id" json:"nimb_id,omitempty"`
	Fname        string       `bson:"fname" json:"fname,omitempty"`
	Lname        string       `bson:"lname" json:"lname,omitempty"`
	Email        string       `bson:"email" json:"email,omitempty"`
	Password     string       `bson:"password" json:"password,omitempty"`
	Organization Organization `bson:"organization" json:"organization,omitempty"`
	Role         Role         `bson:"role" json:"role,omitempty"`
	JWTToken     string       `bson:"-" json:"token,omitempty"`
	Audit        Audit        `bson:"audit" json:"audit"`
}

type Organization struct {
	NIMB_ID string `bson:"nimb_id" json:"nimb_id,omitempty"`
	Name    string `bson:"name" json:"name,omitempty"`
	Address string `bson:"address" json:"address,omitempty"`
	Audit   Audit  `bson:"audit" json:"audit"`
}

type Role struct {
	Name        string   `bson:"name" json:"name,omitempty"`
	Permissions []string `bson:"permissions" json:"permissions,omitempty"`
}

func NewOrganization(c *gin.Context, name, address string) Organization {
	org := Organization{
		NIMB_ID: utils.GenerateNIMBID("N_ORG"),
		Name:    name,
		Address: address,
		Audit: Audit{Active: "ACTIVE",
			IsArchived: false,
			Version:    1, MinorVersion: 0,
			IsProdCandidate: false},
	}
	org.Audit.SetInitialAudit(c)
	return org
}

func NewUser(c *gin.Context, fname, lname, email string, organization string, role Role) User {
	u := User{
		NIMB_ID:      utils.GenerateNIMBID("N_USER"),
		Fname:        fname,
		Lname:        lname,
		Email:        email,
		Organization: Organization{Name: organization},
		Role:         role,
		Audit: Audit{Active: "ACTIVE",
			IsArchived: false,
			Version:    1, MinorVersion: 0,
			IsProdCandidate: false},
	}
	u.Audit.SetInitialAudit(c)
	return u
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
