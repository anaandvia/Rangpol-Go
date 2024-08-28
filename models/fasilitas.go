package models

type Fasilitas struct {
	IdFasilitas     int    `json:"id_fasilitas" gorm:"primaryKey;autoIncrement"`
	NamaFasilitas   string `json:"nama_fasilitas" gorm:"type:varchar(100);not null"`
	DetailFasilitas string `json:"detail_fasilitas" gorm:"type:varchar(100);not null"`
}
