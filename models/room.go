package models

type Room struct {
	Id_room    uint       `json:"id_room" gorm:"primaryKey"`
	No_room    string     `json:"no_room" gorm:"not null;column:no_room;size:50"`
	Name_room  string     `json:"name_room" gorm:"not null;column:name_room;size:50"`
	Lantai     uint       `json:"lantai" gorm:"not null;column:lantai"`
	Kapasitas  uint       `json:"kapasitas" gorm:"not null;column:kapasitas"`
	Foto       string     `json:"foto" gorm:"not null;column:foto;size:255"`
	Status     bool       `json:"status" gorm:"not null;column:status"`
	DetailRoom DetailRoom `gorm:"foreignKey:id_room"`
}
