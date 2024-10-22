package models

type Room struct {
	Id_room    uint         `json:"id_room" gorm:"primaryKey"`
	No_room    string       `json:"no_room" gorm:"not null;column:no_room;size:50"`
	Name_room  string       `json:"name_room" gorm:"not null;column:name_room;size:50"`
	Lantai     uint         `json:"lantai" gorm:"not null;column:lantai"`
	Kapasitas  uint         `json:"kapasitas" gorm:"not null;column:kapasitas"`
	Foto       string       `json:"foto" gorm:"not null;column:foto;size:255"`
	Status     bool         `json:"status" gorm:"not null;column:status"`
	Dlt        int          `gorm:"default:0"`
	DetailRoom DetailRoom   `gorm:"foreignKey:Id_room;references:Id_room"` // Make sure to reference the correct keys
	Peminjaman []Peminjaman `gorm:"foreignKey:IdRoom;references:Id_room"`
}
