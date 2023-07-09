package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	pb "github.com/WiggidyW/weve-esi/proto"
)

func corporationExchangeContractItems(
	c *Client,
	ctx context.Context,
	entity_id uint64,
	contract_id int,
	auth string,
) ([]*pb.ExchangeContractItem, error) {
	return c.exchangeContractItems(
		ctx,
		url.CorporationsCorporationIdContractsContractIdItems,
		entity_id,
		contract_id,
		auth,
	)
}

func characterExchangeContractItems(
	c *Client,
	ctx context.Context,
	entity_id uint64,
	contract_id int,
	auth string,
) ([]*pb.ExchangeContractItem, error) {
	return c.exchangeContractItems(
		ctx,
		url.CharactersCharacterIdContractsContractIdItems,
		entity_id,
		contract_id,
		auth,
	)
}

func (c *Client) exchangeContractItems(
	ctx context.Context,
	url_getter func(uint64, int) string,
	entity_id uint64,
	contract_id int,
	auth string,
) ([]*pb.ExchangeContractItem, error) {
	rep, err := c.crudeRequest(
		ctx,
		url_getter(entity_id, contract_id),
		http.MethodGet,
		auth,
	)
	if err != nil {
		return nil, err
	}
	ec_items := make([]*pb.ExchangeContractItem, 0, len(rep.Json))
	for _, json_ec_item := range rep.Json {
		ec_items = append(ec_items, exchangeContractItemFromJson(json_ec_item))
	}
	return ec_items, nil
}

func exchangeContractItemFromJson(
	json_ec_item map[string]interface{},
) *pb.ExchangeContractItem {
	return &pb.ExchangeContractItem{
		TypeId:   uint32(getValueOrPanic[float64](json_ec_item, "type_id")),
		Quantity: int64(getValueOrPanic[float64](json_ec_item, "quantity")),
	}
}
