package models

type FasilitasRoom struct {
	IdFasilitasRoom int  `json:"id_fasilitasroom" gorm:"primaryKey;autoIncrement"`
	IdFasilitas     int  `json:"id_fasilitas" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdFasilitas;references:IdFasilitas"`
	IdRoom          int  `json:"id_room" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:IdRoom;references:Id_room"`
	Jumlah          int  `json:"jumlah" gorm:"not null"`
	StatusFR        bool `json:"status_fr" gorm:"not null;default:0"`
}
