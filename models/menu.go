package models

type Menu struct {
	Id_menu    uint   `json:"id_menu" gorm:"primaryKey"`
	Parent     string `json:"parent" gorm:"not null;column:parent;size:255"`
	Menu       string `json:"menu" gorm:"not null;column:menu;size:255"`
	Url        string `json:"url" gorm:"not null;column:url;size:255"`
	Jenis      uint   `json:"jenis" gorm:"not null;column:jenis;"`
	Urutan     uint   `json:"urutan" gorm:"not null;column:urutan"`
	Ishide     bool   `json:"ishide" gorm:"not null;default:0"`
	Icon_class string `json:"icon_class" gorm:"not null;column:icon_class;size:255"`
}
