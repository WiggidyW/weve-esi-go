package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

const INCLUDE_COMPLETED = true

type partialIndustryJob struct {
	BlueprintId int64
	IndustryJob *proto.IndustryJob
}

func (c *Client) IndustryJobs(
	ctx context.Context,
	req *proto.IndustryJobsReq,
) (*proto.IndustryJobsRep, error) {
	num_entities := len(req.Corporations) + len(req.Characters)
	chn := make(chan Result[*proto.IndustryJob])

	for _, corporation := range req.Corporations {
		go c.corporationIndustryJobs(
			ctx,
			corporation.Id,
			corporation.Token,
			chn,
		)
	}
	for _, character := range req.Characters {
		go c.characterIndustryJobs(
			ctx,
			character.Id,
			character.Token,
			chn,
		)
	}

	return_rep := new(proto.IndustryJobsRep)
	for num_entities > 0 {
		result := <-chn
		job, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if job != nil {
			return_rep.Inner = append(return_rep.Inner, job)
		} else {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) characterIndustryJobs(
	ctx context.Context,
	character_id uint64,
	token string,
	chn chan Result[*proto.IndustryJob],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}

	blueprint_chn := make(chan Result[blueprints])
	go func(
		ctx context.Context,
		character_id uint64,
		auth string,
		chn chan Result[blueprints],
	) {
		bps, err := c.characterBlueprints(ctx, character_id, auth)
		if err != nil {
			chn <- ResultErr[blueprints](err)
		} else {
			chn <- ResultOk(bps)
		}
	}(ctx, character_id, auth, blueprint_chn)

	jobs_rep, err := c.crudeRequest(
		ctx,
		url.CharactersCharacterIdIndustryJobs(
			character_id,
			INCLUDE_COMPLETED,
		),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}

	blueprints_result := <-blueprint_chn
	blueprints_rep, err := blueprints_result.Unwrap()
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}

	for _, json_job := range jobs_rep.Json {
		partial_job := partialIndustryJobFromJson(json_job)
		job := partial_job.IndustryJob
		blueprint := blueprints_rep[partial_job.BlueprintId]

		if blueprint != nil {
			job.MaterialEfficiency = blueprint.MaterialEfficiency
			job.TimeEfficiency = blueprint.TimeEfficiency
			if blueprint.Runs > 0 {
				job.IsBpc = true
			} else {
				job.IsBpc = false
			}
		}

		chn <- ResultOk(job)
	}

	chn <- ResultNull[*proto.IndustryJob]()
}

func (c *Client) corporationIndustryJobs(
	ctx context.Context,
	corporation_id uint64,
	token string,
	chn chan Result[*proto.IndustryJob],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}

	blueprint_chn := make(chan Result[blueprints])
	go func(
		ctx context.Context,
		character_id uint64,
		auth string,
		chn chan Result[blueprints],
	) {
		bps, err := c.corporationBlueprints(ctx, character_id, auth)
		if err != nil {
			chn <- ResultErr[blueprints](err)
		} else {
			chn <- ResultOk(bps)
		}
	}(ctx, corporation_id, auth, blueprint_chn)

	pages_rep, err := c.crudeRequestHead(
		ctx,
		url.CorporationsCorporationIdIndustryJobs(
			corporation_id,
			INCLUDE_COMPLETED,
			1,
		),
		auth,
	)
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}
	pages := pages_rep.GetPages()

	pages_chn := make(chan Result[*partialIndustryJob])
	for page := 1; page <= pages; page++ {
		go c.corporationIndustryJobsPage(
			ctx,
			corporation_id,
			page,
			auth,
			pages_chn,
		)
	}

	blueprints_result := <-blueprint_chn
	blueprints_rep, err := blueprints_result.Unwrap()
	if err != nil {
		chn <- ResultErr[*proto.IndustryJob](err)
		return
	}

	for pages > 0 {
		result := <-pages_chn
		partial_job, err := result.Unwrap()
		if err != nil {
			chn <- ResultErr[*proto.IndustryJob](err)
			return
		} else if partial_job != nil {
			job := partial_job.IndustryJob
			bp := blueprints_rep[partial_job.BlueprintId]
			if bp != nil {
				job.MaterialEfficiency = bp.MaterialEfficiency
				job.TimeEfficiency = bp.TimeEfficiency
				if bp.Runs > 0 {
					job.IsBpc = true
				} else {
					job.IsBpc = false
				}
			}
			chn <- ResultOk(job)
		} else {
			pages--
		}
	}

	chn <- ResultNull[*proto.IndustryJob]()
}

func (c *Client) corporationIndustryJobsPage(
	ctx context.Context,
	character_id uint64,
	page int,
	auth string,
	chn chan Result[*partialIndustryJob],
) {
	jobs_rep, err := c.crudeRequest(
		ctx,
		url.CorporationsCorporationIdIndustryJobs(
			character_id,
			INCLUDE_COMPLETED,
			page,
		),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*partialIndustryJob](err)
		return
	}

	for _, json_job := range jobs_rep.Json {
		chn <- ResultOk(partialIndustryJobFromJson(json_job))
	}

	chn <- ResultNull[*partialIndustryJob]()
}

func partialIndustryJobFromJson(
	json_job map[string]interface{},
) *partialIndustryJob {
	blueprint_id := int64(getValueOrPanic[float64](json_job, "blueprint_id"))
	industry_job := &proto.IndustryJob{
		LocationId:  uint64(getValueOrPanic[float64](json_job, "facility_id")),
		CharacterId: uint64(getValueOrPanic[float64](json_job, "installer_id")),
		Start:       getTimestampOrPanic(json_job, "start_date"),
		Finish:      getTimestampOrPanic(json_job, "end_date"),
		Probability: getValueOrDefault(json_job, "probability", 1.0),
		ProductId:   uint32(getValueOrDefault(json_job, "product_type_id", 0)),
		BlueprintId: uint32(getValueOrPanic[float64](json_job, "blueprint_type_id")), // blueprint_type_id
		Activity:    int32(getValueOrPanic[float64](json_job, "activity_id")),
		Runs:        int32(getValueOrPanic[float64](json_job, "runs")),
		// MaterialEfficiency: 0,
		// TimeEfficiency: 0,
		// IsBpc: false,
	}
	if industry_job.ProductId == 0 {
		industry_job.ProductId = industry_job.BlueprintId
	}
	return &partialIndustryJob{
		BlueprintId: blueprint_id,
		IndustryJob: industry_job,
	}
}
