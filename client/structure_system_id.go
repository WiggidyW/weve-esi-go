package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
)

func (c *Client) structureSystemId(
	ctx context.Context,
	structure_id uint64,
	auth string,
) (uint32, error) {
	rep, err := c.crudeRequestNoArray(
		ctx,
		url.UniverseStructuresStructureId(structure_id),
		http.MethodGet,
		auth,
	)
	if err != nil {
		return 0, err
	} else if len(rep.Json) == 0 {
		return 0, fmt.Errorf("no json returned from structure id request")
	} else {
		return uint32(getValueOrPanic[float64](rep.Json[0], "solar_system_id")), nil
	}

}
