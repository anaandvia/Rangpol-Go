package models

type User struct {
	Id_user     uint         `json:"id_user" gorm:"primaryKey"`
	Nim         string       `json:"nim" gorm:"not null;column:nim;size:50"`
	Name_user   string       `json:"name_user" gorm:"not null;column:name_user;size:100"`
	Email       string       `json:"email" gorm:"not null;column:email;size:100"`
	Username    string       `json:"username" gorm:"not null;column:username;size:50"`
	Password    string       `json:"password" gorm:"not null"`
	Foto_user   string       `json:"foto_user" gorm:"not null;column:foto_user;size:255"`
	Verif_email bool         `json:"verif_email" gorm:"not null;column:verif_email"`
	Code        string       `json:"code" gorm:"not null"`
	Level       uint         `json:"level" gorm:"not null;column:level"`
	Grup        uint         `json:"grup" gorm:"not null;column:grup"`
	Peminjaman  []Peminjaman `gorm:"foreignKey:IdUser;references:Id_user"` // relasi satu-ke-banyak
}
