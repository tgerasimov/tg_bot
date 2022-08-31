package user_list

import (
	"bot_tg/internal/datastruct"
	"sync"
)

type UserList interface {
	GetUser(id int) (datastruct.User, bool)
	InsertUser(user datastruct.User)
	GetRandomUser() (randomUser datastruct.User, exists bool)
	GetAllUsers() []datastruct.User
}

type userList struct {
	userVault *sync.Map
}

func NewUserList() UserList {
	return &userList{userVault: new(sync.Map)}
}

func (u *userList) GetUser(id int) (datastruct.User, bool) {
	value, exists := u.userVault.Load(id)
	if exists {
		if user, ok := value.(datastruct.User); ok {
			return user, true
		}
	}
	return datastruct.User{}, false
}

func (u *userList) InsertUser(user datastruct.User) {
	u.userVault.Store(user.UserID, user)
}

func (u *userList) GetAllUsers() []datastruct.User {
	res := make([]datastruct.User, 0)
	u.userVault.Range(func(key, value any) bool {
		if us, ok := value.(datastruct.User); ok {
			res = append(res, us)
			return true
		}
		return false
	})
	return res
}

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
