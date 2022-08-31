package repository

import (
	"bot_tg/internal/datastruct"
	"github.com/Masterminds/squirrel"
)

type UserRepo interface {
	Exists(userID int) (bool, error)
	Add(user datastruct.User) error
	UpdateRating(user datastruct.User) error
	Get(userID int) (datastruct.User, error)
}

type userRepo struct {
}

func (u *userRepo) Exists(userID int) (bool, error) {
	_ = pgQb().
		Select("*").
		From(datastruct.UserTableName).
		Where(
			squirrel.Eq{datastruct.UserID_db: userID},
		)
	return true, nil
}

func (u *userRepo) Add(user datastruct.User) error {
	_ = pgQb().
		Insert(datastruct.UserTableName).
		Columns(
			datastruct.UserID_db,
			datastruct.Role_db,
			datastruct.UserName_db,
			datastruct.Rating_db,
		).
		Values(
			user.UserID,
			user.Role,
			user.UserName,
			user.Rating,
		)
	return nil
}

func (u *userRepo) UpdateRating(user datastruct.User) error {
	_ = pgQb().Update(datastruct.UserTableName).
		Where(squirrel.Eq{datastruct.UserID_db: datastruct.User{}}).
		SetMap(map[string]interface{}{
			datastruct.Rating_db: user.Rating,
		})
	return nil
}

func (u *userRepo) Get(userID int) (datastruct.User, error) {
	_ = pgQb().
		Select("*").
		From(datastruct.UserTableName).
		Where(
			squirrel.Eq{datastruct.UserID_db: userID},
		)
	var user datastruct.User

	return user, nil
}
