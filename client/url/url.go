package url

import (
	"fmt"
)

const (
	BASE_URL   = "https://esi.evetech.net/latest"
	DATASOURCE = "tranquility"
)

var (
	BOOL_TO_ORDER_TYPE = map[bool]string{
		false: "sell",
		true:  "buy",
	}
	BOOL_TO_INCLUDE_COMPLETED = map[bool]string{
		false: "false",
		true:  "true",
	}
)

// sort the functions below alphabetically

func IndustrySystems() string {
	return fmt.Sprintf(
		"%s/industry/systems/?datasource=%s",
		BASE_URL,
		DATASOURCE,
	)
}

func MarketsPrices() string {
	return fmt.Sprintf(
		"%s/markets/prices/?datasource=%s",
		BASE_URL,
		DATASOURCE,
	)
}

func MarketsStructuresStructureIdOrders(
	structure_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/markets/structures/%d/?datasource=%s&page=%d",
		BASE_URL,
		structure_id,
		DATASOURCE,
		page,
	)
}

func MarketsRegionIdOrders(
	region_id int,
	page int,
	type_id uint32,
	buy bool,
) string {
	return fmt.Sprintf(
		"%s/markets/%d/orders/?datasource=%s&page=%d&type_id=%d&order_type=%s",
		BASE_URL,
		region_id,
		DATASOURCE,
		page,
		type_id,
		BOOL_TO_ORDER_TYPE[buy],
	)
}

func CharactersCharacterIdAssets(
	character_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/assets/?datasource=%s&page=%d",
		BASE_URL,
		character_id,
		DATASOURCE,
		page,
	)
}

func CorporationsCorporationIdAssets(
	corporation_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/assets/?datasource=%s&page=%d",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		page,
	)
}

func CharactersCharacterIdIndustryJobs(
	corporation_id uint64,
	include_completed bool,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/industry/jobs/?datasource=%s&include_completed=%s",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		BOOL_TO_INCLUDE_COMPLETED[include_completed],
	)
}

func CorporationsCorporationIdIndustryJobs(
	corporation_id uint64,
	include_completed bool,
	page int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/industry/jobs/?datasource=%s&include_completed=%s&page=%d",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		BOOL_TO_INCLUDE_COMPLETED[include_completed],
		page,
	)
}

func CharactersCharacterIdOrders(
	character_id uint64,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/orders/?datasource=%s",
		BASE_URL,
		character_id,
		DATASOURCE,
	)
}

func CorporationsCorporationIdOrders(
	corporation_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/orders/?datasource=%s&page=%d",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		page,
	)
}

func CharactersCharacterIdBlueprints(
	character_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/blueprints/?datasource=%s&page=%d",
		BASE_URL,
		character_id,
		DATASOURCE,
		page,
	)
}

func CorporationsCorporationIdBlueprints(
	corporation_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/blueprints/?datasource=%s&page=%d",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		page,
	)
}

func CharactersCharacterIdWalletTransactions(
	character_id uint64,
	from_id int64,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/wallet/transactions/?datasource=%s&from_id=%d",
		BASE_URL,
		character_id,
		DATASOURCE,
		from_id,
	)
}

func CorporationsCorporationIdWalletsDivisionTransactions(
	corporation_id uint64,
	from_id int64,
	division int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/wallets/%d/transactions/?datasource=%s&from_id=%d",
		BASE_URL,
		corporation_id,
		division,
		DATASOURCE,
		from_id,
	)
}

func CharactersCharacterIdSkills(
	character_id uint64,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/skills/?datasource=%s",
		BASE_URL,
		character_id,
		DATASOURCE,
	)
}

func CorporationsCorporationIdContracts(
	corporation_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/contracts/?datasource=%s&page=%d",
		BASE_URL,
		corporation_id,
		DATASOURCE,
		page,
	)
}

func CharactersCharacterIdContracts(
	character_id uint64,
	page int,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/contracts/?datasource=%s&page=%d",
		BASE_URL,
		character_id,
		DATASOURCE,
		page,
	)
}

func CorporationsCorporationIdContractsContractIdItems(
	corporation_id uint64,
	contract_id int,
) string {
	return fmt.Sprintf(
		"%s/corporations/%d/contracts/%d/items/?datasource=%s",
		BASE_URL,
		corporation_id,
		contract_id,
		DATASOURCE,
	)
}

func CharactersCharacterIdContractsContractIdItems(
	character_id uint64,
	contract_id int,
) string {
	return fmt.Sprintf(
		"%s/characters/%d/contracts/%d/items/?datasource=%s",
		BASE_URL,
		character_id,
		contract_id,
		DATASOURCE,
	)
}

// func UniverseStationsStationId(
// 	station_id uint64,
// ) string {
// 	panic("unimpl")
// }

// func UniverseSystemsSystemId(
// 	system_id int,
// ) string {
// 	panic("unimpl")
// }

// func UniverseConstellationsConstellationId(
// 	constellation_id int,
// ) string {
// 	panic("unimpl")
// }
