package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
)

type blueprint struct {
	BlueprintId        int64
	Runs               int32
	MaterialEfficiency int32
	TimeEfficiency     int32
}

type blueprints map[int64]*blueprint

func (c *Client) corporationBlueprints(
	ctx context.Context,
	corporation_id uint64,
	auth string,
) (blueprints, error) {
	return c.blueprints(
		ctx,
		url.CorporationsCorporationIdBlueprints,
		corporation_id,
		auth,
	)
}

func (c *Client) characterBlueprints(
	ctx context.Context,
	character_id uint64,
	auth string,
) (blueprints, error) {
	return c.blueprints(
		ctx,
		url.CharactersCharacterIdBlueprints,
		character_id,
		auth,
	)
}

func (c *Client) blueprints(
	ctx context.Context,
	url_getter func(uint64, int) string,
	entity_id uint64,
	auth string,
) (blueprints, error) {
	pages_rep, err := c.crudeRequestHead(
		ctx,
		url_getter(entity_id, 1),
		auth,
	)
	if err != nil {
		return nil, err
	}
	pages := pages_rep.GetPages()

	chn := make(chan Result[*blueprint])
	for page := 1; page <= pages; page++ {
		go c.blueprintsPage(
			ctx,
			url_getter(entity_id, page),
			auth,
			chn,
		)
	}

	blueprints_rep := make(blueprints)
	for pages > 0 {
		result := <-chn
		blueprint, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if blueprint != nil {
			blueprints_rep[blueprint.BlueprintId] = blueprint
		} else {
			pages--
		}
	}

	return blueprints_rep, nil
}

func (c *Client) blueprintsPage(
	ctx context.Context,
	url string,
	auth string,
	chn chan Result[*blueprint],
) {
	bps_rep, err := c.crudeRequest(ctx, url, http.MethodGet, auth)
	if err != nil {
		chn <- ResultErr[*blueprint](err)
		return
	}

	for _, json_bp := range bps_rep.Json {
		chn <- ResultOk(blueprintFromJson(json_bp))
	}

	chn <- ResultNull[*blueprint]()
}

func blueprintFromJson(
	json_blueprint map[string]interface{},
) *blueprint {
	return &blueprint{
		BlueprintId:        int64(getValueOrPanic[float64](json_blueprint, "item_id")),
		Runs:               int32(getValueOrPanic[float64](json_blueprint, "runs")),
		MaterialEfficiency: int32(getValueOrPanic[float64](json_blueprint, "material_efficiency")),
		TimeEfficiency:     int32(getValueOrPanic[float64](json_blueprint, "time_efficiency")),
	}
}
