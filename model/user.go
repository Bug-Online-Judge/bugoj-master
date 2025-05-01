package model

import "time"

type Role struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `gorm:"uniqueIndex;not null" json:"name"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `json:"-"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Roles     []Role    `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

type UserDTO struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func ToUserDTO(u User) UserDTO {
	roleNames := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		roleNames[i] = r.Name
	}
	return UserDTO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Roles:    roleNames,
	}
}
