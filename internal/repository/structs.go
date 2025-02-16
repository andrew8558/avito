package repository

type User struct {
	Login    string `db:"login"`
	Password string `db:"password"`
	Coins    int32  `db:"coins"`
}
