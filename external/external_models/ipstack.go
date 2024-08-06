package external_models

type IPStackResolveIPResponse struct {
	Ip            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

type IPStackResolveIPResponseLocation struct {
	GeonameId               int                                        `json:"geoname_id"`
	Capital                 string                                     `json:"capital"`
	Languages               []IPStackResolveIPResponseLocationLanguage `json:"languages"`
	CountryFlag             string                                     `json:"country_flag"`
	CountryFlagEmoji        string                                     `json:"country_flag_emoji"`
	CountryFlagEmojiUnicode string                                     `json:"country_flag_emoji_unicode"`
	CallingCode             string                                     `json:"calling_code"`
	IsEu                    bool                                       `json:"is_eu"`
}

type IPStackResolveIPResponseLocationLanguage struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Native string `json:"native"`
}
