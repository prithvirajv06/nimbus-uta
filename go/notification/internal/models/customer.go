package models

type User struct {
	ID           string       `bson:"_id,omitempty" json:"id"`
	Fname        string       `bson:"fname" json:"fname"`
	Lname        string       `bson:"lname" json:"lname"`
	Email        string       `bson:"email" json:"email"`
	Password     string       `bson:"password" json:"-"`
	Organization Organization `bson:"organization" json:"organization"`
	Role         Role         `bson:"role" json:"role"`
	Active       bool         `bson:"active" json:"active"`
	Audit        Audit        `bson:"audit" json:"audit"`
}

type Organization struct {
	ID      string `bson:"_id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
	Active  bool   `bson:"active" json:"active"`
	Audit   Audit  `bson:"audit" json:"audit"`
}

type Role struct {
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

func (r Role) CreateDefaultAdminRole() Role {
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

func NewUser(fname, lname, email, password, organization string, role Role) *User {
	return &User{
		Fname:        fname,
		Lname:        lname,
		Email:        email,
		Password:     password,
		Organization: Organization{Name: organization},
		Role:         role,
		Active:       true,
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
