package dao

type User interface {
	GetUser()
	DeleteUser()
	AddUser()
	UpdateUser()
}

type user struct {
}

func (u user) GetUser() {
	//TODO implement me
	panic("implement me")
}

func (u user) DeleteUser() {
	//TODO implement me
	panic("implement me")
}

func (u user) AddUser() {
	//TODO implement me
	panic("implement me")
}

func (u user) UpdateUser() {
	//TODO implement me
	panic("implement me")
}
