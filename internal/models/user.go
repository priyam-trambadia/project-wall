package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) Insert() {
	err := DB.Ping()
	if err != nil {
		fmt.Println("Errrrrrrr")
	}
	fmt.Println("doneeeeeeee")
}
