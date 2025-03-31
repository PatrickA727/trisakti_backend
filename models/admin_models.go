package models

type AdminUser struct {
	ID			int		`json:"id"`
	Username	string	`json:"username" binding:"required"`
	Password	string	`json:"password" binding:"required"`
}

type Sessions struct {
	AdminID 		int		`json:"student_id" binding:"required"`
	RefreshToken 	string		`json:"refresh_token" binding:"required"`
}
