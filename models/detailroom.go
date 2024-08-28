package models

type DetailRoom struct {
	Id_detailroom uint   `json:"id_detailroom" gorm:"primaryKey"`
	Id_room       uint   `json:"id_room" gorm:"not null;column:id_room;unique;constrait:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Fungsi        string `json:"fungsi" gorm:"not null;column:fungsi;size:100"`
	Keperluan     string `json:"keperluan" gorm:"not null;column:keperluan;size:50"`
	PIC           string `json:"pic" gorm:"not null;column:pic;size:100"`
	Koorlab       string `json:"koorlab" gorm:"not null;column:koorlab;size:100"`
}
