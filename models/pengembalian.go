package models

import "time"

type Pengembalian struct {
	Id              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	TglPengembalian time.Time `json:"tgl_pengembalian" gorm:"type:datetime"`
	IdPeminjaman    uint      `json:"id_peminjaman" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdPeminjaman;references:IdPeminjaman"`
	FotoB           string    `json:"foto_b" gorm:"type:varchar(50);not null"`
	Kendala         string    `json:"kendala" gorm:"type:text;not null"`
	StatusKembali   bool      `json:"status_kembali" gorm:"not null;default:0"`
	Alasan          string    `json:"alasan" gorm:"type:text;not null"`
	TglAccBack      time.Time `json:"tgl_acc_back" gorm:"type:date"`
}
