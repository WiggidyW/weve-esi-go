package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

type skill struct {
	CharacterId uint64
	SkillId     uint32
	Level       uint32
}

func (c *Client) Skills(
	ctx context.Context,
	req *proto.SkillsReq,
) (*proto.SkillsRep, error) {
	num_entities := len(req.Characters)
	chn := make(chan Result[*skill])

	for _, character := range req.Characters {
		go c.characterSkills(
			ctx,
			character.Id,
			character.Token,
			chn,
		)
	}

	return_rep := new(proto.SkillsRep)
	for num_entities > 0 {
		result := <-chn
		skill, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if skill != nil {
			return_rep.Inner[skill.CharacterId].
				Inner[skill.SkillId] = skill.Level
		} else {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) characterSkills(
	ctx context.Context,
	character_id uint64,
	token string,
	chn chan Result[*skill],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*skill](err)
		return
	}

	skills_rep, err := c.crudeRequest(
		ctx,
		url.CharactersCharacterIdSkills(character_id),
		http.MethodGet,
		auth,
	)
	if err != nil {
		chn <- ResultErr[*skill](err)
		return
	}

	for _, json_skill := range skills_rep.Json {
		chn <- ResultOk(skillFromJson(json_skill, character_id))
	}

	chn <- ResultNull[*skill]()
}

func skillFromJson(
	json_skill map[string]interface{},
	character_id uint64,
) *skill {
	return &skill{
		CharacterId: character_id,
		SkillId:     uint32(getValueOrPanic[float64](json_skill, "skill_id")),
		Level:       uint32(getValueOrPanic[float64](json_skill, "active_skill_level")),
	}
}
