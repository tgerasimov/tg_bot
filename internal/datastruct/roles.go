package datastruct

const (
	Default = iota
	Admin
)

type Role struct {
	roleName        string
	roleDescription string
}

var Roles = map[int]Role{
	Default: {roleName: "default", roleDescription: "default role"},
	Admin:   {roleName: "admin", roleDescription: "admin role"},
}
