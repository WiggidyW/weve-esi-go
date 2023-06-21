package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

func (c *Client) AdjustedPrice(
	ctx context.Context,
	req *proto.AdjustedPriceReq,
) (*proto.AdjustedPriceRep, error) {
	items_rep, err := c.crudeRequest(
		ctx,
		url.MarketsPrices(),
		http.MethodGet,
		NULL_AUTH,
	)
	if err != nil {
		return nil, err
	}

	return_rep := new(proto.AdjustedPriceRep)
	for _, item := range items_rep.Json {
		type_id, adj_price := adjustedPriceFromJson(item)
		return_rep.Inner[type_id] = adj_price
	}

	return return_rep, nil
}

func adjustedPriceFromJson(
	item map[string]interface{},
) (uint32, float64) {
	type_id := uint32(getValueOrPanic[float64](item, "type_id"))
	adj_price := getValueOrPanic[float64](item, "adjusted_price")
	return type_id, adj_price
}
