package appTypes

type RoleID int

const (
	Guest RoleID = iota // 游客，自增
	User                // 普通用户
	Admin               // 管理员
)
