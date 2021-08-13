package fair

type Model struct {
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	Setcens        string  `json:"setcens"`
	Areap          string  `json:"areap"`
	CodDistrict    string  `json:"cod_district"`
	District       string  `gorm:"index" json:"district"`
	CodSubCityHall string  `json:"cod_sub_city_hall"`
	SubCityHall    string  `json:"sub_city_hall"`
	Region5        string  `gorm:"index" json:"region_5"`
	Region8        string  `json:"region_8"`
	Name           string  `gorm:"index" json:"name"`
	Registry       string  `gorm:"uniqueIndex" json:"registry"`
	Address        string  `json:"address"`
	AddressNumber  string  `json:"address_number"`
	Neighborhood   string  `gorm:"index" json:"neighborhood"`
	Landmark       string  `json:"landmark"`
}

func (Model) TableName() string {
	return "streetfair"
}
