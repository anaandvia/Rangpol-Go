package models

import "time"

type Peminjaman struct {
	IdPeminjaman  int       `json:"id_peminjaman" gorm:"primaryKey;autoIncrement"`
	IdUser        int       `json:"id_user" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdUser;references:Id"`
	IdRoom        int       `json:"id_room" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdRoom;references:Id_room"`
	NamaKegiatan  string    `json:"nama_kegiatan" gorm:"type:varchar(100);not null"`
	TglAcara      time.Time `json:"tgl_acara" gorm:"not null"`
	TglAkhirAcara time.Time `json:"tgl_akhir_acara" gorm:"not null"`
	Status        bool      `json:"status" gorm:"not null;default:0"`
	TglAcc        time.Time `json:"tgl_acc"`
}
