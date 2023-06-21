package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

type locationAndFlag struct {
	LocationId uint64
	Flag       string
}

type extendedAsset struct {
	LocationId uint64
	ItemId     uint64
	TypeId     uint32
	Asset      *proto.Asset
}

func (e *extendedAsset) toLocationAndFlag() *locationAndFlag {
	return &locationAndFlag{
		LocationId: e.LocationId,
		Flag:       e.Asset.Flags[0],
	}
}

func (c *Client) Assets(
	ctx context.Context,
	req *proto.AssetsReq,
) (*proto.AssetsRep, error) {
	num_entities := len(req.Corporations) + len(req.Characters)
	chn := make(chan Result[*extendedAsset])

	for _, corporation := range req.Corporations {
		go c.entityAssets(
			ctx,
			url.CorporationsCorporationIdAssets,
			url.CharactersCharacterIdBlueprints,
			corporation.Id,
			corporation.Token,
			chn,
		)
	}
	for _, character := range req.Characters {
		go c.entityAssets(
			ctx,
			url.CharactersCharacterIdAssets,
			url.CharactersCharacterIdBlueprints,
			character.Id,
			character.Token,
			chn,
		)
	}

	return_rep := new(proto.AssetsRep)
	for num_entities > 0 {
		result := <-chn
		asset, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if asset != nil {
			location_assets := return_rep.Inner[asset.LocationId]
			type_assets := location_assets.Inner[asset.TypeId]
			type_assets.Inner = append(
				type_assets.Inner,
				asset.Asset,
			)
		} else {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) entityAssets(
	ctx context.Context,
	url_getter func(uint64, int) string,
	bp_url_getter func(uint64, int) string,
	entity_id uint64,
	token string,
	chn chan Result[*extendedAsset],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*extendedAsset](err)
		return
	}

	blueprint_chn := make(chan Result[blueprints])
	go func(
		ctx context.Context,
		url_getter func(uint64, int) string,
		entity_id uint64,
		auth string,
		chn chan Result[blueprints],
	) {
		bps, err := c.blueprints(ctx, url_getter, entity_id, auth)
		if err != nil {
			chn <- ResultErr[blueprints](err)
		} else {
			chn <- ResultOk(bps)
		}
	}(ctx, bp_url_getter, entity_id, auth, blueprint_chn)

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url_getter(entity_id, 1),
		auth,
	)
	if err != nil {
		chn <- ResultErr[*extendedAsset](err)
		return
	}
	pages := pages_rep.GetPages()

	page_chn := make(chan Result[*extendedAsset])
	for page := 1; page <= pages; page++ {
		go c.assetsPage(
			ctx,
			url_getter(entity_id, page),
			entity_id,
			auth,
			page_chn,
		)
	}

	blueprints_result := <-blueprint_chn
	blueprints_rep, err := blueprints_result.Unwrap()
	if err != nil {
		chn <- ResultErr[*extendedAsset](err)
		return
	}

	assets := make([]*extendedAsset, pages*1000)
	asset_flags := make(map[uint64]*locationAndFlag, pages*1000)
	for pages > 0 {
		result := <-page_chn
		asset, err := result.Unwrap()
		if err != nil {
			chn <- ResultErr[*extendedAsset](err)
			return
		} else if asset != nil {
			asset_flags[asset.ItemId] = asset.toLocationAndFlag()
			assets = append(assets, asset)
			blueprint := blueprints_rep[int64(asset.ItemId)]
			if blueprint != nil {
				extendWithBlueprint(asset.Asset, blueprint)
			}
		} else {
			pages--
		}
	}

	for _, asset := range assets {
		for {
			parent := asset_flags[asset.LocationId]
			if parent == nil {
				break
			}
			asset.LocationId = parent.LocationId
			asset.Asset.Flags = append(
				asset.Asset.Flags,
				parent.Flag,
			)
		}
		chn <- ResultOk(asset)
	}

	chn <- ResultNull[*extendedAsset]()
}

func (c *Client) assetsPage(
	ctx context.Context,
	url string,
	entity_id uint64,
	auth string,
	chn chan Result[*extendedAsset],
) {
	assets_rep, err := c.crudeRequest(ctx, url, http.MethodGet, auth)
	if err != nil {
		chn <- ResultErr[*extendedAsset](err)
		return
	}
	for _, json_asset := range assets_rep.Json {
		chn <- ResultOk(extendedAssetFromJson(json_asset, entity_id))
	}
	chn <- ResultNull[*extendedAsset]()
}

func extendedAssetFromJson(
	json_asset map[string]interface{},
	entity_id uint64,
) *extendedAsset {
	location_id := uint64(getValueOrPanic[float64](json_asset, "location_id"))
	item_id := uint64(getValueOrPanic[float64](json_asset, "item_id"))
	type_id := uint32(getValueOrPanic[float64](json_asset, "type_id"))
	asset := &proto.Asset{
		EntityId:           entity_id,
		Quantity:           int64(getValueOrPanic[float64](json_asset, "quantity")),
		Runs:               0,
		MaterialEfficiency: 0,
		TimeEfficiency:     0,
		Flags:              []string{getValueOrPanic[string](json_asset, "location_flag")},
	}
	return &extendedAsset{
		Asset:      asset,
		LocationId: location_id,
		ItemId:     item_id,
		TypeId:     type_id,
	}
}

func extendWithBlueprint(
	asset *proto.Asset,
	blueprint *blueprint,
) {
	asset.Runs = blueprint.Runs
	asset.MaterialEfficiency = blueprint.MaterialEfficiency
	asset.TimeEfficiency = blueprint.TimeEfficiency
}
