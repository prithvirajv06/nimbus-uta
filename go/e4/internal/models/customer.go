package models

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
