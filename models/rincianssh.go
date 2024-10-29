package models

import "gorm.io/gorm"

type Rincianssh struct {
	RinciansshID   string  `json:"rinciansshid"`
	SSrinciansshID string  `json:"ssrinciansshid"`
	RinciansshKode string  `json:"rinciansshkode"`
	RinciansshNama string  `json:"rinciansshnama"`
	Spesifikasi    string  `json:"spesifikasi"`
	Satuan         string  `json:"satuan"`
	Harga          string  `json:"harga"`
	Thang          string  `json:"thang"`
	OpAdd          string  `json:"opadd"`
	PcAdd          string  `json:"pcadd"`
	TglAdd         string  `json:"tgladd"`
	OpEdit         string  `json:"opedit"`
	PcEdit         string  `json:"pcedit"`
	TglEdit        string  `json:"tgledit"`
	Dlt            bool    `json:"dlt"`
	Merk           *string `json:"merk"` // Gunakan pointer jika bisa null
}

func GetPaginatedRincianssh(db *gorm.DB, page int, pageSize int) ([]Rincianssh, int64, error) {
	var rinciansshList []Rincianssh
	var total int64

	offset := (page - 1) * pageSize

	// Hitung total data untuk pagination
	db.Model(&Rincianssh{}).Count(&total)

	// Ambil data dengan limit dan offset untuk pagination
	err := db.Limit(pageSize).Offset(offset).Find(&rinciansshList).Error
	if err != nil {
		return nil, 0, err
	}

	return rinciansshList, total, nil
}
