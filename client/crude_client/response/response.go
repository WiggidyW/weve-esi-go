package response

import (
	"time"
)

type EsiHeadResponse struct {
	Pages   int
	Expires time.Time
}

func (r *EsiHeadResponse) GetPages() int {
	return r.Pages
}

type EsiResponse struct {
	Json []map[string]interface{}
	// Json    interface{}
	Etag    string
	Expires time.Time
}

// func (r *EsiResponse) StationSystemId() int {
// 	return int(r.Json.(map[string]interface{})["system_id"].(float64))
// }

// func (r *EsiResponse) SystemConstellationId() int {
// 	return int(
// 		r.Json.(map[string]interface{})["constellation_id"].(float64),
// 	)
// }

// func (r *EsiResponse) ConstellationRegionId() int {
// 	return int(r.Json.(map[string]interface{})["region_id"].(float64))
// }

// func (r *EsiResponse) MarketOrders()
