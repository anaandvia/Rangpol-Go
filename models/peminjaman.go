package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Peminjaman struct {
	IdPeminjaman           uint             `json:"id_peminjaman" gorm:"primaryKey;autoIncrement"`
	IdUser                 uint             `json:"id_user" gorm:"not null;index"`
	IdRoom                 uint             `json:"id_room" gorm:"not null;index"`
	NamaKegiatan           string           `json:"nama_kegiatan" gorm:"type:varchar(100);not null"`
	TglAcara               time.Time        `json:"tgl_acara" gorm:"not null"`
	TglAkhirAcara          time.Time        `json:"tgl_akhir_acara" gorm:"not null"`
	Status                 uint             `json:"status" gorm:"not null;default:0"`
	TglAcc                 time.Time        `json:"tgl_acc"`
	DetailPeminjaman       DetailPeminjaman `gorm:"foreignKey:IdPeminjaman;references:IdPeminjaman"` // relasi satu-ke-satu
	User                   User             `gorm:"foreignKey:IdUser;references:Id_user"`
	Room                   Room             `gorm:"foreignKey:IdRoom;references:Id_room"`
	TglAcaraFormatted      string
	TglAkhirAcaraFormatted string
	TglAcaraDay            string         `json:"tgl_acara_day"`
	TglAkhirAcaraDay       string         `json:"tgl_akhir_acara_day"`
	Pengembalian           []Pengembalian `gorm:"foreignKey:IdPeminjaman"`
}

func (p *Peminjaman) BeforeSave(tx *gorm.DB) (err error) {
	if p.TglAkhirAcara.Before(p.TglAcara) {
		return errors.New("TglAkhirAcara cannot be before TglAcara")
	}
	return nil
}
