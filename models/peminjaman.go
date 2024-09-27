package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Peminjaman struct {
	IdPeminjaman           uint              `json:"id_peminjaman" gorm:"primaryKey;autoIncrement"`
	IdUser                 uint              `json:"id_user" gorm:"not null;index"`
	IdRoom                 uint              `json:"id_room" gorm:"not null;index"`
	NamaKegiatan           string            `json:"nama_kegiatan" gorm:"type:varchar(100);not null"`
	TglAcara               time.Time         `json:"tgl_acara" gorm:"not null"`
	TglAkhirAcara          time.Time         `json:"tgl_akhir_acara" gorm:"not null"`
	Status                 uint              `json:"status" gorm:"not null;default:0"`
	TglAcc                 time.Time         `json:"tgl_acc"`
	DetailPeminjaman       *DetailPeminjaman `gorm:"foreignKey:IdPeminjaman;references:IdPeminjaman"` // relasi satu-ke-satu
	User                   *User             `gorm:"foreignKey:IdUser;references:Id_user"`
	Room                   *Room             `gorm:"foreignKey:IdRoom;references:Id_room"`
	TglAcaraFormatted      string
	TglAkhirAcaraFormatted string
	TglAcaraDay            string        `json:"tgl_acara_day"`
	TglAkhirAcaraDay       string        `json:"tgl_akhir_acara_day"`
	Pengembalian           *Pengembalian `gorm:"foreignKey:IdPeminjaman"`
	Dlt                    int           `gorm:"default:0"`
	KetVerif               string        `json:"ketverif" gorm:"type:text;"`
	KeteranganTelat        string
}

func (p *Peminjaman) BeforeSave(tx *gorm.DB) (err error) {
	// Mengecek apakah operasi ini adalah "insert" (penambahan data)
	isInsert := p.IdPeminjaman == 0
	// Jika operasi adalah penambahan data (insert), lakukan validasi
	if isInsert {
		// Validasi: TglAkhirAcara tidak boleh sebelum TglAcara
		if p.TglAkhirAcara.Before(p.TglAcara) {
			return errors.New("TglAkhirAcara cannot be before TglAcara")
		}

		// Validasi: Cek apakah ada peminjaman lain yang tumpang tindih
		var count int64
		err = tx.Model(&Peminjaman{}).
			Where("id_room = ? AND ((tgl_acara <= ? AND tgl_akhir_acara >= ?) OR (tgl_acara <= ? AND tgl_akhir_acara >= ?)) AND status != ?",
				p.IdRoom, p.TglAkhirAcara, p.TglAcara, p.TglAkhirAcara, p.TglAcara, 2).
			Count(&count).Error
		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("Room is already booked for the specified time period")
		}
	}

	return nil
}
