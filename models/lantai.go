package models

type Lantai struct {
	Id_lantai uint `json:"id_lantai" gorm:"primaryKey"`
	No_lantai uint `json:"no_lantai" gorm:"not null;column:no_lantai"`
}
