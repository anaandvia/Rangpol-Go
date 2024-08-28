package models

type Jurusan struct {
	Id_jurusan   uint   `json:"id_jurusan" gorm:"primaryKey"`
	Kode_Jurusan string `json:"kode_jurusan" gorm:"not null;column:kode_jurusan;size:20"`
	Nama_Jurusan string `json:"nama_jurusan" gorm:"not null;column:nama_jurusan;size:255"`
}
