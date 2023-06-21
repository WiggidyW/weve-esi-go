package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

func (c *Client) stationMarketOrders(
	ctx context.Context,
	req *proto.MarketOrdersReq,
) (*proto.MarketOrdersRep, error) {
	region_id, err := c.dbGetRegionId(ctx, req.LocationId)
	if err != nil {
		return nil, err
	}

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url.MarketsRegionIdOrders(region_id, 1, req.TypeId, req.Buy),
		NULL_AUTH,
	)
	if err != nil {
		return nil, err
	}
	pages := pages_rep.GetPages()

	chn := make(chan Result[*proto.MarketOrder])
	for page := 1; page <= pages; page++ {
		go c.stationMarketOrdersPage(ctx, region_id, page, req, chn)
	}

	return_rep := new(proto.MarketOrdersRep)
	for pages > 0 {
		result := <-chn
		order, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if order != nil {
			return_rep.Inner = append(return_rep.Inner, order)
		} else {
			pages--
		}
	}

	return return_rep, nil
}

func (c *Client) stationMarketOrdersPage(
	ctx context.Context,
	region_id int,
	page int,
	req *proto.MarketOrdersReq,
	chn chan Result[*proto.MarketOrder],
) {
	orders_rep, err := c.crudeRequest(
		ctx,
		url.MarketsRegionIdOrders(region_id, page, req.TypeId, req.Buy),
		http.MethodGet,
		NULL_AUTH,
	)
	if err != nil {
		chn <- ResultErr[*proto.MarketOrder](err)
		return
	}

	for _, order := range orders_rep.Json {
		location_id, market_order := stationMarketOrderFromJson(order)
		if location_id == req.LocationId {
			chn <- ResultOk(market_order)
		}
	}

	chn <- ResultNull[*proto.MarketOrder]()
}

func stationMarketOrderFromJson(
	order map[string]interface{},
) (uint64, *proto.MarketOrder) {
	location_id := uint64(getValueOrPanic[float64](order, "location_id"))
	market_order := &proto.MarketOrder{
		Quantity: int64(getValueOrPanic[float64](order, "volume_remain")),
		Price:    getValueOrPanic[float64](order, "price"),
	}
	return location_id, market_order
}
