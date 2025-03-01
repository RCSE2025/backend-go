package model

//- id			PK[int]
//- INN		Unique[int]
//OPTIONAL:
//- OGRN		Unique[BigInt]
//- owner		str
//- short_name 	str # Название орги в кратком виде: ООО "лалалала"
//- full_name	str # Название орги в полном виде: общество с ограниченной ответственностью "лалалалала"
//- address		json

type Business struct {
	baseModel
	ID        int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	INN       int64   `json:"inn" gorm:"unique;not null"`
	OGRN      *int64  `json:"ogrn,omitempty" gorm:"unique"`
	Owner     *string `json:"owner,omitempty" gorm:"size:100"`
	ShortName *string `json:"short_name,omitempty" gorm:"size:100"`
	FullName  *string `json:"full_name,omitempty" gorm:"size:100"`
	Address   *string `json:"address,omitempty" gorm:"size:10000"`
}

func (b *Business) TableName() string {
	return "businesses"
}
