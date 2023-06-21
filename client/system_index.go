package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

func (c *Client) SystemIndex(
	ctx context.Context,
	req *proto.SystemIndexReq,
) (*proto.SystemIndexRep, error) {
	indices_rep, err := c.crudeRequest(
		ctx,
		url.IndustrySystems(),
		http.MethodGet,
		NULL_AUTH,
	)
	if err != nil {
		return nil, err
	}

	return_rep := new(proto.SystemIndexRep)
	for _, system := range indices_rep.Json {
		system_id, system_index := systemIndexFromJson(system)
		return_rep.Inner[system_id] = system_index
	}

	return return_rep, nil
}

func systemIndexFromJson(
	system map[string]interface{},
) (uint32, *proto.SystemIndex) {
	system_id := uint32(getValueOrPanic[float64](system, "solar_system_id"))
	cost_indices := getValueOrPanic[[]map[string]interface{}](system, "cost_indices")

	system_index := new(proto.SystemIndex)
	for _, cost_index := range cost_indices {
		index := getValueOrPanic[float64](cost_index, "cost_index")
		activity := getValueOrPanic[string](cost_index, "activity")
		switch activity {
		case "manufacturing":
			system_index.Manufacturing = index
		case "researching_time_efficiency":
			system_index.ResearchTe = index
		case "researching_material_efficiency":
			system_index.ResearchMe = index
		case "copying":
			system_index.Copying = index
		case "invention":
			system_index.Invention = index
		case "reaction":
			system_index.Reactions = index
		default: // Do nothing
		}
	}

	return system_id, system_index
}
