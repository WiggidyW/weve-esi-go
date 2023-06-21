package client

import (
	"context"
	"net/http"

	"github.com/WiggidyW/weve-esi/client/crude_client/response"
	"github.com/WiggidyW/weve-esi/client/url"
	"github.com/WiggidyW/weve-esi/proto"
)

var WALLET_DIVISIONS = [7]int{1, 2, 3, 4, 5, 6, 7}

type extTransactions struct {
	LastId       int64 // 0 if no more transactions to fetch
	Transactions []*extTransaction
}

type extTransaction struct {
	LocationId  uint64
	TypeId      uint32
	Transaction *proto.Transaction
}

func (c *Client) Transactions(
	ctx context.Context,
	req *proto.TransactionsReq,
) (*proto.TransactionsRep, error) {
	num_entities := len(req.Characters) + len(req.Corporations)*7
	chn := make(chan Result[*extTransaction])

	for _, corporation := range req.Corporations {
		for _, division := range WALLET_DIVISIONS {
			go c.corporationTransactions(
				ctx,
				corporation.Id,
				division,
				corporation.Token,
				req.Since,
				0,
				chn,
			)
		}
	}
	for _, character := range req.Characters {
		go c.characterTransactions(
			ctx,
			character.Id,
			character.Token,
			req.Since,
			0,
			chn,
		)
	}

	return_rep := new(proto.TransactionsRep)
	for num_entities > 0 {
		result := <-chn
		t, err := result.Unwrap()
		if err != nil {
			return nil, err
		} else if t != nil {
			loctn_transactions := return_rep.Inner[t.LocationId]
			type_transactions := loctn_transactions.Inner[t.TypeId]
			type_transactions.Inner = append(
				type_transactions.Inner,
				t.Transaction,
			)
		} else {
			num_entities--
		}
	}

	return return_rep, nil
}

func (c *Client) characterTransactions(
	ctx context.Context,
	character_id uint64,
	token string,
	since uint64,
	from_id int64,
	chn chan Result[*extTransaction],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*extTransaction](err)
		return
	}

	transactions, err := c.entityTransactionsFromId(
		ctx,
		url.CharactersCharacterIdWalletTransactions(
			character_id,
			from_id,
		),
		auth,
		since,
		from_id == 0,
	)
	if err != nil {
		chn <- ResultErr[*extTransaction](err)
		return
	}

	var child_chn chan Result[*extTransaction] = nil
	if transactions.LastId != 0 {
		child_chn = make(chan Result[*extTransaction])
		go c.characterTransactions(
			ctx,
			character_id,
			auth,
			since,
			transactions.LastId,
			child_chn,
		)
	}

	c.entityTransactions(transactions, chn, child_chn)
}

func (c *Client) corporationTransactions(
	ctx context.Context,
	corporation_id uint64,
	division int,
	token string,
	since uint64,
	from_id int64,
	chn chan Result[*extTransaction],
) {
	auth, err := c.crudeRequestAuth(ctx, token)
	if err != nil {
		chn <- ResultErr[*extTransaction](err)
		return
	}

	transactions, err := c.entityTransactionsFromId(
		ctx,
		url.CorporationsCorporationIdWalletsDivisionTransactions(
			corporation_id,
			from_id,
			division,
		),
		auth,
		since,
		from_id == 0,
	)
	if err != nil {
		chn <- ResultErr[*extTransaction](err)
		return
	}

	var child_chn chan Result[*extTransaction] = nil
	if transactions.LastId != 0 {
		child_chn = make(chan Result[*extTransaction])
		go c.corporationTransactions(
			ctx,
			corporation_id,
			division,
			auth,
			since,
			transactions.LastId,
			child_chn,
		)
	}

	c.entityTransactions(transactions, chn, child_chn)
}

func (c *Client) entityTransactionsFromId(
	ctx context.Context,
	url string,
	auth string,
	since uint64,
	use_cache bool,
) (*extTransactions, error) {
	var trns_rep *response.EsiResponse
	var err error

	if use_cache {
		trns_rep, err = c.crudeRequest(ctx, url, http.MethodGet, auth)
	} else {
		trns_rep, err = c.
			crudeRequestNoCache(ctx, url, http.MethodGet, auth)
	}
	if err != nil {
		return nil, err
	}

	return extTransactionsFromJson(trns_rep.Json, since), nil
}

func (c *Client) entityTransactions(
	transactions *extTransactions,
	chn chan Result[*extTransaction],
	child_chn chan Result[*extTransaction],
) {
	for _, transaction := range transactions.Transactions {
		chn <- ResultOk(transaction)
	}

	if child_chn != nil {
		for {
			result := <-child_chn
			transaction, err := result.Unwrap()
			if err != nil {
				chn <- ResultErr[*extTransaction](err)
				return
			} else if transaction != nil {
				chn <- ResultOk(transaction)
			} else {
				break
			}
		}
	}

	chn <- ResultNull[*extTransaction]()
}

// binary search function to find the index of the first transaction that is older than the given timestamp
func findFirstOlderTransaction(
	transactions []map[string]interface{},
	since uint64,
) int {
	// ) (int, bool) { // returns index, found
	var start_pt = 0
	var end_pt = len(transactions) - 1

	if timestampFromJson(transactions[end_pt]) > since {
		return len(transactions) // no transactions older than since
	}
	if timestampFromJson(transactions[start_pt]) <= since {
		return start_pt // all transactions older than since
	}

	for start_pt != end_pt {
		mid_pt := (start_pt + end_pt) / 2
		if mid_pt == start_pt {
			start_pt = end_pt
		} else if timestampFromJson(transactions[mid_pt]) > since {
			start_pt = mid_pt
		} else { // timestamp <= since
			end_pt = mid_pt
		}
	}

	return end_pt
	// return start_pt, true
}

func extTransactionsFromJson(
	transactions []map[string]interface{},
	since uint64,
) *extTransactions {
	first_older_idx := findFirstOlderTransaction(transactions, since)
	ext_transactions := &extTransactions{
		Transactions: make([]*extTransaction, first_older_idx),
	}

	for i := 0; i < first_older_idx; i++ {
		ext_transactions.Transactions = append(
			ext_transactions.Transactions,
			extTransactionFromJson(transactions[i]),
		)
		if i == len(transactions)-1 {
			ext_transactions.LastId = transactionIdFromJson(
				transactions[i],
			)
		}
	}

	return ext_transactions
}

func extTransactionFromJson(transaction map[string]interface{}) *extTransaction {
	return &extTransaction{
		LocationId: uint64(getValueOrPanic[float64](transaction, "location_id")),
		TypeId:     uint32(getValueOrPanic[float64](transaction, "type_id")),
		Transaction: &proto.Transaction{
			Quantity: int64(getValueOrPanic[float64](transaction, "quantity")),
			Price:    getValueOrPanic[float64](transaction, "unit_price"),
			Buy:      getValueOrPanic[bool](transaction, "is_buy"),
		},
	}
}

func timestampFromJson(transaction map[string]interface{}) uint64 {
	return getTimestampOrPanic(transaction, "date")
}

func transactionIdFromJson(transaction map[string]interface{}) int64 {
	return int64(getValueOrPanic[float64](transaction, "transaction_id"))
}
