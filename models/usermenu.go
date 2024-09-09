package models

type Usermenu struct {
	IdUserMenu uint `json:"id_usermenu" gorm:"primaryKey"`
	Id_level   uint `json:"id_level"`
	Id_menu    uint `json:"id_menu" gorm:"not null;index"`
	View       uint `json:"view"`
	Ubah       uint `json:"ubah"`
	Hapus      uint `json:"hapus"`
}
