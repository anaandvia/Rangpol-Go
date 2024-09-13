package models

type Usermenu struct {
	IdUserMenu uint `json:"id_usermenu" gorm:"primaryKey"`
	Id_menu    uint `json:"id_menu" gorm:"not null;index"`
	Id_akses   uint `json:"id_akses"`
	View       uint `json:"view"`
	Add        uint `json:"add"`
	Edit       uint `json:"edit"`
	Delete     uint `json:"delete"`
	Print      uint `json:"print"`
}
