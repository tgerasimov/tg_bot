package datastruct

const UserTableName = "user_table"

const (
	UserID_db   = "user_id"
	UserName_db = "user_name"
	Role_db     = "user_role"
	Rating_db   = "user_rating"
)

type User struct {
	UserID   int
	UserName string
	Role     Role
	Rating   float64
}
