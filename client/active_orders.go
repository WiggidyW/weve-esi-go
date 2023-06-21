package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

type extActiveOrder struct {
	LocationId  uint64
	TypeId      uint32
	ActiveOrder *proto.ActiveOrder
}

func (c *Client) ActiveOrders(
	ctx context.Context,
	req *proto.ActiveOrdersReq,
) (*proto.ActiveOrdersRep, error) {
	num_entities := len(req.Characters) + len(req.Corporations)
	chn := make(chan Result[*extActiveOrder])

	for _, corporation := range req.Corporations {
		go c.corporationActiveOrders(
			ctx,
			corporation.Id,
			corporation.Token,
			chn,
		)
	}
	for _, character := range req.Characters {
		go c.characterActiveOrders(
			ctx,
			character.Id,
			character.Token,
			chn,
		)
	}

	return_rep := new(proto.ActiveOrdersRep)
	for num_entities > 0 {
		result := <-chn
		order, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if order != nil {
			location_orders := return_rep.Inner[order.LocationId]
			type_orders := location_orders.Inner[order.TypeId]
			type_orders.Inner = append(
				type_orders.Inner,
				order.ActiveOrder,
			)
		} else {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) characterActiveOrders(
	ctx context.Context,
	character_id uint64,
	token string,
	chn chan Result[*extActiveOrder],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*extActiveOrder](err)
		return
	}

	orders_rep, err := c.crudeRequest(
		ctx,
		url.CharactersCharacterIdOrders(character_id),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*extActiveOrder](err)
		return
	}

	for _, order := range orders_rep.Json {
		chn <- ResultOk(activeOrderFromJson(order))
	}

	chn <- ResultNull[*extActiveOrder]()
}

func (c *Client) corporationActiveOrders(
	ctx context.Context,
	corporation_id uint64,
	token string,
	chn chan Result[*extActiveOrder],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*extActiveOrder](err)
		return
	}

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url.CorporationsCorporationIdOrders(corporation_id, 1),
		auth,
	)
	if err != nil {
		chn <- ResultErr[*extActiveOrder](err)
		return
	}
	pages := pages_rep.GetPages()

	page_chn := make(chan Result[*extActiveOrder])
	for page := 1; page <= pages; page++ {
		go c.corporationActiveOrdersPage(
			ctx,
			corporation_id,
			page,
			auth,
			page_chn,
		)
	}

	for pages > 0 {
		result := <-page_chn
		order, err := result.Unwrap()
		if err != nil {
			chn <- result
			return
		} else if order != nil {
			chn <- result
		} else {
			pages--
		}
	}

	chn <- ResultNull[*extActiveOrder]()
}

func (c *Client) corporationActiveOrdersPage(
	ctx context.Context,
	corporation_id uint64,
	page int,
	auth string,
	chn chan Result[*extActiveOrder],
) {
	orders_rep, err := c.crudeRequest(
		ctx,
		url.CorporationsCorporationIdOrders(corporation_id, page),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*extActiveOrder](err)
		return
	}

	for _, order := range orders_rep.Json {
		chn <- ResultOk(activeOrderFromJson(order))
	}

	chn <- ResultNull[*extActiveOrder]()
}

func activeOrderFromJson(order map[string]interface{}) *extActiveOrder {
	return &extActiveOrder{
		LocationId: uint64(getValueOrPanic[float64](order, "location_id")),
		TypeId:     uint32(getValueOrPanic[float64](order, "type_id")),
		ActiveOrder: &proto.ActiveOrder{
			// LocationId: uint64(getValueOrPanic[float64](order, "location_id")),
			// TypeId:     uint32(getValueOrPanic[float64](order, "type_id")),
			Quantity: int64(getValueOrPanic[float64](order, "volume_remain")),
			Price:    getValueOrPanic[float64](order, "price"),
			Buy:      getValueOrDefault(order, "is_buy_order", false),
		},
	}
}
