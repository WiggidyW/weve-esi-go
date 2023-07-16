package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/crude_client"
	"github.com/WiggidyW/weve-esi/client/url"
	pb "github.com/WiggidyW/weve-esi/proto"
)

const EXCHANGE_CONTRACT_TYPE = "item_exchange"
const CONTRACT_ACTIVE_STATUS = "outstanding"

type partialExchangeContract struct {
	ContractId       int
	ExchangeContract *pb.ExchangeContract
}

type exchangeContractLocationData struct {
	SystemId uint32
	RegionId uint32
}

type itemGetter = func(
	c *Client,
	ctx context.Context,
	entity_id uint64,
	contract_id int,
	auth string,
) ([]*pb.ExchangeContractItem, error)

func (c *Client) ExchangeContracts(
	ctx context.Context,
	req *pb.ExchangeContractsReq,
) (*pb.ExchangeContractsRep, error) {
	num_entities := len(req.Corporations) + len(req.Characters)
	chn := make(chan Result[*pb.ExchangeContract])

	for _, corporation := range req.Corporations {
		go c.corporationExchangeContracts(
			ctx,
			corporation.Id,
			corporation.Token,
			req.StructureToken,
			req.IncludeItems,
			req.ActiveOnly,
			chn,
		)
	}
	for _, character := range req.Characters {
		go c.characterExchangeContracts(
			ctx,
			character.Id,
			character.Token,
			req.StructureToken,
			req.IncludeItems,
			req.ActiveOnly,
			chn,
		)
	}

	return_rep := new(pb.ExchangeContractsRep)
	for num_entities > 0 {
		result := <-chn
		contract, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if contract != nil {
			return_rep.Inner = append(return_rep.Inner, contract)
		} else /* if err == nil && contract == nil */ {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) corporationExchangeContracts(
	ctx context.Context,
	corporation_id uint64,
	token string,
	structure_token string,
	include_items bool,
	active_only bool,
	chn chan Result[*pb.ExchangeContract],
) {
	c.exchangeContracts(
		ctx,
		url.CorporationsCorporationIdContracts,
		corporationExchangeContractItems,
		include_items,
		active_only,
		corporation_id,
		token,
		structure_token,
		chn,
	)
}

func (c *Client) characterExchangeContracts(
	ctx context.Context,
	character_id uint64,
	token string,
	structure_token string,
	include_items bool,
	active_only bool,
	chn chan Result[*pb.ExchangeContract],
) {
	c.exchangeContracts(
		ctx,
		url.CharactersCharacterIdContracts,
		characterExchangeContractItems,
		include_items,
		active_only,
		character_id,
		token,
		structure_token,
		chn,
	)
}

func (c *Client) exchangeContracts(
	ctx context.Context,
	url_getter func(uint64, int) string,
	item_getter itemGetter,
	include_items bool,
	active_only bool,
	entity_id uint64,
	token string,
	structure_token string,
	chn chan Result[*pb.ExchangeContract],
) {
	var structure_auth string
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*pb.ExchangeContract](err)
		return
	}
	if structure_token == "" {
		structure_auth = auth
	} else {
		structure_auth, err = c.crudeRequestAuth(ctx, structure_token)
		if err != nil {
			chn <- ResultErr[*pb.ExchangeContract](err)
			return
		}
	}

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url_getter(entity_id, 1),
		auth,
	)
	if err != nil {
		chn <- ResultErr[*pb.ExchangeContract](err)
		return
	}
	pages := pages_rep.GetPages()

	page_chn := make(chan Result[*pb.ExchangeContract])
	for page := 1; page <= pages; page++ {
		go c.exchangeContractsPage(
			ctx,
			url_getter,
			item_getter,
			include_items,
			active_only,
			entity_id,
			page_chn,
			page,
			auth,
			structure_auth,
		)
	}

	for pages > 0 {
		result := <-page_chn
		contract, err := result.Unwrap()
		if err != nil {
			chn <- ResultErr[*pb.ExchangeContract](err)
			return
		} else if contract != nil {
			chn <- ResultOk[*pb.ExchangeContract](contract)
		} else /* if err == nil && contract == nil */ {
			pages--
		}
	}

	chn <- ResultNull[*pb.ExchangeContract]()
}

func (c *Client) exchangeContractsPage(
	ctx context.Context,
	url_getter func(uint64, int) string,
	item_getter itemGetter,
	include_items bool,
	active_only bool,
	entity_id uint64,
	chn chan Result[*pb.ExchangeContract],
	page int,
	auth string,
	structure_auth string,
) {
	ec_contracts_rep, err := c.crudeRequest(
		ctx,
		url_getter(entity_id, page),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*pb.ExchangeContract](err)
		return
	}

	items_chn := make(chan Result[*pb.ExchangeContract])
	num_contracts := 0

	for _, json_contract := range ec_contracts_rep.Json {
		// Filter out non-exchange and (if active_only) non-active contracts
		contract_type := getValueOrPanic[string](json_contract, "type")
		if contract_type != EXCHANGE_CONTRACT_TYPE {
			continue
		}
		if active_only {
			status := getValueOrPanic[string](json_contract, "status")
			if status != CONTRACT_ACTIVE_STATUS {
				continue
			}
		}
		// Send out a request for the items in the contract
		num_contracts++
		partial_contract := partialExchangeContractFromJson(json_contract)
		go c.exchangeContractsContract(
			ctx,
			item_getter,
			include_items,
			partial_contract,
			entity_id,
			auth,
			structure_auth,
			items_chn,
		)
	}

	// Wait for all the items to come back
	for num_contracts > 0 {
		result := <-items_chn
		contract, err := result.Unwrap()
		if err != nil {
			chn <- ResultErr[*pb.ExchangeContract](err)
			return
		} else /* if contract != nil */ {
			num_contracts--
			chn <- ResultOk[*pb.ExchangeContract](contract)
		}
	}

	chn <- ResultNull[*pb.ExchangeContract]()
}

func (c *Client) exchangeContractsContract(
	ctx context.Context,
	item_getter itemGetter,
	include_items bool,
	partial_contract *partialExchangeContract,
	entity_id uint64,
	auth string,
	structure_auth string,
	chn chan Result[*pb.ExchangeContract],
) {
	location_chn := make(chan Result[*exchangeContractLocationData])
	go c.exchangeContractsLocationData(
		ctx,
		partial_contract.ExchangeContract.LocationId,
		structure_auth,
		location_chn,
	)

	if include_items {
		item_chn := make(chan Result[*[]*pb.ExchangeContractItem])
		go func(
			c *Client,
			ctx context.Context,
			item_getter itemGetter,
			entity_id uint64,
			contract_id int,
			auth string,
			item_chn chan Result[*[]*pb.ExchangeContractItem],
		) {
			items, err := item_getter(c, ctx, entity_id, contract_id, auth)
			if err != nil {
				item_chn <- ResultErr[*[]*pb.ExchangeContractItem](err)
			} else {
				item_chn <- ResultOk(&items)
			}
		}(
			c,
			ctx,
			item_getter,
			entity_id,
			partial_contract.ContractId,
			auth,
			item_chn,
		)
		items_rep := <-item_chn
		items, err := items_rep.Unwrap()
		if err != nil {
			chn <- ResultErr[*pb.ExchangeContract](err)
			return
		} else {
			partial_contract.ExchangeContract.Items = *items
		}
	}

	location_rep := <-location_chn
	location_data, err := location_rep.Unwrap()
	if err != nil {
		chn <- ResultErr[*pb.ExchangeContract](err)
		return
	} else {
		partial_contract.ExchangeContract.SystemId = location_data.SystemId
		partial_contract.ExchangeContract.RegionId = location_data.RegionId
	}

	chn <- ResultOk[*pb.ExchangeContract](partial_contract.ExchangeContract)
}

func (c *Client) exchangeContractsLocationData(
	ctx context.Context,
	location_id uint64,
	auth string,
	chn chan Result[*exchangeContractLocationData],
) {
	// Get the system ID
	first_digit := u64FirstDigit(location_id)
	var system_id uint32
	var err error
	if first_digit == 1 {
		// Get it from ESI as it's a structure
		system_id, err = c.structureSystemId(ctx, location_id, auth)
		// Check if it's a status error
		if err != nil {
			esi_err, ok := err.(crude_client.StatusError)
			if ok && esi_err.Code == http.StatusForbidden {
				// If it is, return 0,0
				chn <- ResultOk(&exchangeContractLocationData{
					SystemId: 0,
					RegionId: 0,
				})
			}
		}

	} else if first_digit == 6 {
		// Get it from static DB as it's a station
		system_id, err = c.dbGetStationSystemId(ctx, location_id)
	}
	if err != nil {
		chn <- ResultErr[*exchangeContractLocationData](err)
		return
	}

	// Get the region ID
	region_id, err := c.dbGetSystemRegionId(ctx, system_id)
	if err != nil {
		chn <- ResultErr[*exchangeContractLocationData](err)
		return
	}

	chn <- ResultOk(&exchangeContractLocationData{
		SystemId: system_id,
		RegionId: region_id,
	})
}

func partialExchangeContractFromJson(
	json_ec_contract map[string]interface{},
) *partialExchangeContract {
	contract_id := int(getValueOrPanic[float64](json_ec_contract, "contract_id"))
	exchange_contract := &pb.ExchangeContract{
		LocationId:  uint64(getValueOrPanic[float64](json_ec_contract, "start_location_id")),
		Description: getValueOrPanic[string](json_ec_contract, "title"),
		Price:       getValueOrPanic[float64](json_ec_contract, "price"),
		Reward:      getValueOrPanic[float64](json_ec_contract, "reward"),
		Expires:     getTimestampOrPanic(json_ec_contract, "date_expired"),
		Issued:      getTimestampOrPanic(json_ec_contract, "date_issued"),
		Volume:      getValueOrPanic[float64](json_ec_contract, "volume"),
		CharId:      uint32(getValueOrPanic[float64](json_ec_contract, "issuer_id")),
		CorpId:      uint32(getValueOrPanic[float64](json_ec_contract, "issuer_corporation_id")),
		IsCorp:      getValueOrPanic[bool](json_ec_contract, "for_corporation"),
	}
	return &partialExchangeContract{
		ContractId:       contract_id,
		ExchangeContract: exchange_contract,
	}
}
