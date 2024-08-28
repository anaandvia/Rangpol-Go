package models

type DetailPeminjaman struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	IdPeminjaman int    `json:"id_peminjaman" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdPeminjaman;references:IdPeminjaman"`
	PJ           string `json:"pj" gorm:"type:varchar(100);not null"`
	PA           string `json:"pa" gorm:"type:varchar(100);not null"`
	PK           string `json:"pk" gorm:"type:varchar(100);not null"`
	NTamu        int    `json:"n_tamu" gorm:"not null"`
	SifatAcara   string `json:"sifat_acara" gorm:"type:varchar(50);not null"`
	JenisAcara   string `json:"jenis_acara" gorm:"type:varchar(50);not null"`
	Keterangan   string `json:"keterangan" gorm:"type:text;not null"`
}
