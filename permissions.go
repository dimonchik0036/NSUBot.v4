package main

const (
	PermissionAdmin = 666
	PermissionVIP   = 10
	PermissionUser  = 0
)

func (u *Users) ChangePermission(subscriber string, permission int) bool {
	user := u.User(subscriber)
	if user == nil {
		return false
	}
	user.Permission = permission
	return true
}
