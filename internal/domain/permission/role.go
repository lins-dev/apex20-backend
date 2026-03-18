package permission

type Role string

const (
	RoleGM      Role = "gm"
	RolePlayer  Role = "player"
	RoleTrusted Role = "trusted"
)
