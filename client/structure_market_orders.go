package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

func (c *Client) structureMarketOrders(
	ctx context.Context,
	req *proto.MarketOrdersReq,
) (*proto.MarketOrdersRep, error) {
	auth, err := c.crudeRequestAuth(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url.MarketsStructuresStructureIdOrders(req.LocationId, 1),
		auth,
	)
	if err != nil {
		return nil, err
	}
	pages := pages_rep.GetPages()

	chn := make(chan Result[*proto.MarketOrder])
	for page := 1; page <= pages; page++ {
		go c.structureMarketOrdersPage(ctx, auth, page, req, chn)
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

func (c *Client) structureMarketOrdersPage(
	ctx context.Context,
	auth string,
	page int,
	req *proto.MarketOrdersReq,
	chn chan Result[*proto.MarketOrder],
) {
	orders_rep, err := c.crudeRequest(
		ctx,
		url.MarketsStructuresStructureIdOrders(req.LocationId, page),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*proto.MarketOrder](err)
		return
	}

	for _, order := range orders_rep.Json {
		chn <- ResultOk(structureMarketOrderFromJson(order))
	}

	chn <- ResultNull[*proto.MarketOrder]()
}

func structureMarketOrderFromJson(
	order map[string]interface{},
) *proto.MarketOrder {
	return &proto.MarketOrder{
		Quantity: int64(getValueOrPanic[float64](order, "volume_remain")),
		Price:    getValueOrPanic[float64](order, "price"),
	}
}
