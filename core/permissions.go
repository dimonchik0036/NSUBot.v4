package core

const (
	PermissionAdmin = 666
	PermissionVIP   = 10
	PermissionUser  = 0
)

func (u *Users) ChangePermission(subscriber string, permission int) {
	user := u.User(subscriber)
	user.Permission = permission
}
