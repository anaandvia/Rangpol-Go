package models

type Menu struct {
	Id_menu    uint       `json:"id_menu" gorm:"primaryKey"`
	Id_Parent  uint       `json:"id_parent" gorm:"null;column:id_parent;"`
	Menu       string     `json:"menu" gorm:"not null;column:menu;size:255"`
	Url        string     `json:"url" gorm:"not null;column:url;size:255"`
	Jenis      uint       `json:"jenis" gorm:"not null;column:jenis;"`
	Urutan     uint       `json:"urutan" gorm:"not null;column:urutan"`
	Ishide     bool       `json:"ishide" gorm:"not null;default:0"`
	Icon_class string     `json:"icon_class" gorm:"not null;column:icon_class;size:255"`
	Children   []Menu     `gorm:"foreignKey:Id_Parent;references:Id_menu"`
	UserMenu   []Usermenu `gorm:"foreignKey:Id_menu;references:Id_menu"`
}
