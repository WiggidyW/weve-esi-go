package client

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/proto"
)

func (c *Client) MarketOrders(
	ctx context.Context,
	req *proto.MarketOrdersReq,
) (*proto.MarketOrdersRep, error) {
	first_digit := u64FirstDigit(req.LocationId)
	if first_digit == 1 {
		return c.structureMarketOrders(ctx, req)
	} else if first_digit == 6 {
		return c.stationMarketOrders(ctx, req)
	} else {
		return nil, fmt.Errorf(
			"MarketOrdersReq had invalid LocationId: %d",
			req.LocationId,
		)
	}
}

// func (c *Client) stationIdToRegionId(
// 	ctx context.Context,
// 	station_id uint64,
// ) (int, error) {
// 	var rep *EsiResponse
// 	var err error

// 	rep, err = c.InnerClient.Request(
// 		ctx,
// 		urlUniverseStationsStationId(station_id),
// 		http.MethodGet,
// 		"",
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	system_id := rep.StationSystemId()

// 	rep, err = c.InnerClient.Request(
// 		ctx,
// 		urlUniverseSystemsSystemId(system_id),
// 		http.MethodGet,
// 		"",
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	constellation_id := rep.SystemConstellationId()

// 	rep, err = c.InnerClient.Request(
// 		ctx,
// 		urlUniverseConstellationsConstellationId(constellation_id),
// 		http.MethodGet,
// 		"",
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	region_id := rep.ConstellationRegionId()

// 	return region_id, nil
// }
