package user_list

import (
	"sync"

	"bot_tg/internal/datastruct"
)

type UserList interface {
	GetUser(id int) (datastruct.User, bool)
	InsertUser(user datastruct.User)
	GetRandomUser() (randomUser datastruct.User, exists bool)
	GetAllUsers() map[int]datastruct.User
	Exists(id int) bool
}

type userList struct {
	userVault *sync.Map
}

func NewUserList() UserList {
	return &userList{userVault: new(sync.Map)}
}

//GetUser ...
func (u *userList) GetUser(id int) (datastruct.User, bool) {
	value, exists := u.userVault.Load(id)
	if exists {
		if user, ok := value.(datastruct.User); ok {
			return user, true
		}
	}
	return datastruct.User{}, false
}

//InsertUser ...
func (u *userList) InsertUser(user datastruct.User) {
	u.userVault.Store(user.UserID, user)
}

//GetAllUsers ...
func (u *userList) GetAllUsers() map[int]datastruct.User {
	res := make(map[int]datastruct.User)
	u.userVault.Range(func(key, value any) bool {
		if us, ok := value.(datastruct.User); ok {
			res[us.UserID] = us
			return true
		}
		return false
	})
	return res
}

//GetRandomUser ...
func (u *userList) GetRandomUser() (randomUser datastruct.User, exists bool) {
	u.userVault.Range(func(key, value any) bool {
		if us, ok := value.(datastruct.User); ok {
			randomUser = us
			exists = true
			return true
		}
		exists = false
		return false
	})
	return randomUser, exists
}

//Exists check does user exists
func (u *userList) Exists(id int) bool {
	value, ex := u.userVault.Load(id)
	if ex {
		if _, ok := value.(datastruct.User); ok {
			return true
		}
		return false
	}
	return false
}
