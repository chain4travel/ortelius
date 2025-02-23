// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code whose
// original notices appear below.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************
// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/chain4travel/magellan/caching"
	"github.com/chain4travel/magellan/cfg"
	"github.com/chain4travel/magellan/services/indexes/params"
	"github.com/chain4travel/magellan/utils"
	"github.com/gocraft/web"
)

const (
	DefaultLimit       = 1000
	DefaultOffsetLimit = 10000
)

type V2Context struct {
	*Context
	version uint8
	chainID *ids.ID
}

const (
	MetricCount  = "api_count"
	MetricMillis = "api_millis"
)

const (
	MetricTransactionsCount   = "api_transactions_count"
	MetricTransactionsMillis  = "api_transactions_millis"
	MetricCTransactionsCount  = "api_ctransactions_count"
	MetricCTransactionsMillis = "api_ctransactions_millis"
	MetricAddressesCount      = "api_addresses_count"
	MetricAddressesMillis     = "api_addresses_millis"
	MetricAddressChainsCount  = "api_address_chains_count"
	MetricAddressChainsMillis = "api_address_chains_millis"
	MetricAggregateCount      = "api_aggregate_count"
	MetricAggregateMillis     = "api_aggregate_millis"
	MetricAssetCount          = "api_asset_count"
	MetricAssetMillis         = "api_asset_millis"
	MetricSearchCount         = "api_search_count"
	MetricSearchMillis        = "api_search_millis"
)

// AddV2Routes mounts a V2 API router at the given path, displaying the given
// indexBytes at the root. If chainID is not nil the handlers run in v1
// compatible mode where the `version` param is set to "1" and requests to
// default to filtering by the given chainID.
func AddV2Routes(ctx *Context, router *web.Router, path string, indexBytes []byte, chainID *ids.ID) {
	utils.Prometheus.CounterInit(MetricCount, MetricCount)
	utils.Prometheus.CounterInit(MetricMillis, MetricMillis)

	utils.Prometheus.CounterInit(MetricTransactionsCount, MetricTransactionsCount)
	utils.Prometheus.CounterInit(MetricTransactionsMillis, MetricTransactionsMillis)

	utils.Prometheus.CounterInit(MetricCTransactionsCount, MetricCTransactionsCount)
	utils.Prometheus.CounterInit(MetricTransactionsMillis, MetricCTransactionsMillis)

	utils.Prometheus.CounterInit(MetricAddressesCount, MetricAddressesCount)
	utils.Prometheus.CounterInit(MetricAddressesMillis, MetricAddressesMillis)

	utils.Prometheus.CounterInit(MetricAddressChainsCount, MetricAddressChainsCount)
	utils.Prometheus.CounterInit(MetricAddressChainsMillis, MetricAddressChainsMillis)

	utils.Prometheus.CounterInit(MetricAggregateCount, MetricAggregateCount)
	utils.Prometheus.CounterInit(MetricAggregateMillis, MetricAggregateMillis)

	utils.Prometheus.CounterInit(MetricAssetCount, MetricAssetCount)
	utils.Prometheus.CounterInit(MetricAssetMillis, MetricAssetMillis)

	utils.Prometheus.CounterInit(MetricSearchCount, MetricSearchCount)
	utils.Prometheus.CounterInit(MetricSearchMillis, MetricSearchMillis)

	v2ctx := V2Context{Context: ctx}
	router.Subrouter(v2ctx, path).
		Get("/", func(c *V2Context, resp web.ResponseWriter, _ *web.Request) {
			if _, err := resp.Write(indexBytes); err != nil {
				ctx.sc.Log.Warn("response write failed",
					zap.Error(err),
				)
			}
		}).

		// Handle legacy v1 logic
		Middleware(func(c *V2Context, w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
			c.version = 2
			if chainID != nil {
				c.chainID = chainID
				c.version = 1
			}
			next(w, r)
		}).
		Get("/search", (*V2Context).Search).
		Get("/aggregates", (*V2Context).Aggregate).
		Get("/txfeeAggregates", (*V2Context).TxfeeAggregate).
		Get("/transactions/aggregates", (*V2Context).Aggregate).
		Get("/addressChains", (*V2Context).AddressChains).
		Post("/addressChains", (*V2Context).AddressChainsPost).
		Post("/validatorsInfo", (*V2Context).ValidatorsInfo).
		Get("/activeAddresses", (*V2Context).ActiveAddresses).
		Get("/uniqueAddresses", (*V2Context).UniqueAddresses).
		Get("/averageBlockSize", (*V2Context).AverageBlockSize).
		Get("/dailyTransactions", (*V2Context).DailyTransactions).
		Get("/dailyGasUsed", (*V2Context).DailyGasUsed).
		Get("/avgGasPriceUsed", (*V2Context).AvgGasPriceUsed).
		Get("/dailyTokenTransfer", (*V2Context).DailyTokenTransfer).
		Get("/dailyEmissions", (*V2Context).DailyEmissions).
		Get("/networkEmissions", (*V2Context).NetworkEmissions).
		Get("/transactionEmissions", (*V2Context).TransactionEmissions).
		Get("/countryEmissions", (*V2Context).CountryEmissions).

		// List and Get routes
		Get("/transactions", (*V2Context).ListTransactions).
		Post("/transactions", (*V2Context).ListTransactionsPost).
		Get("/transactions/:id", (*V2Context).GetTransaction).
		Get("/addresses", (*V2Context).ListAddresses).
		Get("/addresses/:id", (*V2Context).GetAddress).
		Get("/outputs", (*V2Context).ListOutputs).
		Get("/outputs/:id", (*V2Context).GetOutput).
		Get("/assets", (*V2Context).ListAssets).
		Get("/assets/:id", (*V2Context).GetAsset).
		Get("/atxdata/:id", (*V2Context).ATxData).
		Get("/ptxdata/:id", (*V2Context).PTxData).
		Get("/ctxdata/:id", (*V2Context).CTxData).
		Get("/cblocks", (*V2Context).ListCBlocks).
		Get("/ctransactions", (*V2Context).ListCTransactions).
		Get("/cacheaddresscounts", (*V2Context).CacheAddressCounts).
		Get("/cachetxscounts", (*V2Context).CacheTxCounts).
		Get("/cacheassets", (*V2Context).CacheAssets).
		Get("/cacheassetaggregates", (*V2Context).CacheAssetAggregates).
		Get("/cacheaggregates/:id", (*V2Context).CacheAggregates).
		Get("/multisigalias/:owners", (*V2Context).GetMultisigAlias).
		Post("/rewards", (*V2Context).GetRewardPost)
}

// AVAX
func (c *V2Context) ValidatorsInfo(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricSearchMillis),
		utils.NewCounterIncCollect(MetricSearchCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ValidatorParams{}
	if err := p.SetParamInfo(c.version, c.sc.ServicesCfg.CaminoNode); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("geoIPValidatorsInfo", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return utils.GetValidatorsGeoIPInfo(p.RPC, &c.sc.Services.GeoIP, c.sc.Logger())
		},
	})
}

func (c *V2Context) ActiveAddresses(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams("activeAddresses", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ActiveAddresses(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) UniqueAddresses(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams("uniqueAddresses", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.UniqueAddresses(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) AverageBlockSize(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams("AverageBlockSize", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.AverageBlockSizeReader(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) DailyTransactions(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("dailyTransactionsChart %s", p.ListParams.ID)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.DailyTransactions(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) DailyGasUsed(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("DailyGasUsed %s", p.ListParams.ID)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GasUsedPerDay(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) AvgGasPriceUsed(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("AvgGasPriceUsed %s", p.ListParams.ID)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.AvgGasPriceUsed(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) DailyTokenTransfer(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.StatisticsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("DailyTokenTransfer %s", p.ListParams.ID)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.DailyTokenTransfer(ctx, &p.ListParams)
		},
	})
}

func (c *V2Context) DailyEmissions(w web.ResponseWriter, r *web.Request) {
	p := &params.EmissionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("Daily Emissions %s", p.ListParams.EndTime)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: time.Duration(c.sc.ServicesCfg.CacheEmissionsInterval) * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return utils.GetDailyEmissions(p.ListParams.StartTime, p.ListParams.EndTime, c.sc.Services.InmutableInsights, c.sc.ServicesCfg.CaminoNode), nil
		},
	})
}

func (c *V2Context) CountryEmissions(w web.ResponseWriter, r *web.Request) {
	p := &params.EmissionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("Country Emissions %s", p.ListParams.EndTime)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: time.Duration(c.sc.ServicesCfg.CacheEmissionsInterval) * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return utils.GetCountryEmissions(p.ListParams.StartTime, p.ListParams.EndTime, c.sc.Services.InmutableInsights, c.sc.ServicesCfg.CaminoNode)
		},
	})
}

func (c *V2Context) NetworkEmissions(w web.ResponseWriter, r *web.Request) {
	p := &params.EmissionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	key := fmt.Sprintf("Network Emissions %s", p.ListParams.EndTime)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: time.Duration(c.sc.ServicesCfg.CacheEmissionsInterval) * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return utils.GetNetworkEmissions(p.ListParams.StartTime, p.ListParams.EndTime, c.sc.Services.InmutableInsights, c.sc.ServicesCfg.CaminoNode)
		},
	})
}

func (c *V2Context) TransactionEmissions(w web.ResponseWriter, r *web.Request) {
	p := &params.EmissionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	key := fmt.Sprintf("Transaction Emissions %s", p.ListParams.EndTime)
	c.WriteCacheable(w, caching.Cacheable{
		TTL: time.Duration(c.sc.ServicesCfg.CacheEmissionsInterval) * time.Hour,
		Key: c.cacheKeyForParams(key, p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return utils.GetNetworkEmissionsPerTransaction(p.ListParams.StartTime, p.ListParams.EndTime, c.sc.Services.InmutableInsights, c.sc.ServicesCfg.CaminoNode)
		},
	})
}

func (c *V2Context) Search(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricSearchMillis),
		utils.NewCounterIncCollect(MetricSearchCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.SearchParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("search", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.Search(ctx, p, c.avaxAssetID)
		},
	})
}

func (c *V2Context) TxfeeAggregate(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.TxfeeAggregateParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)
	if len(p.ChainIDs) == 0 {
		c.WriteErr(w, 400, fmt.Errorf("chainID is required"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("aggregate_txfee", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.TxfeeAggregate(c.sc.AggregatesCache, p)
		},
	})
}

func (c *V2Context) Aggregate(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.AggregateParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)
	if len(p.ChainIDs) == 0 {
		c.WriteErr(w, 400, fmt.Errorf("chainID is required"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("aggregate", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.Aggregate(c.sc.AggregatesCache, p)
		},
	})
}

func (c *V2Context) GetMultisigAlias(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	ownersParam := r.PathParams["owners"]
	ownerAddresses := []string{}
	for _, ownerAddress := range strings.Split(ownersParam, ",") {
		// convert owner address from bech32 to be used internally
		addr, err := address.ParseToID(ownerAddress)
		if err != nil {
			c.WriteErr(w, 400, err)
			return
		}
		ownerAddresses = append(ownerAddresses, addr.String())
	}
	// calculate cache key from owner addresses
	cacheKey := fmt.Sprintf("%x", sha256.Sum256([]byte(ownersParam)))
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForID("multisig_alias", cacheKey),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetMultisigAlias(ctx, ownerAddresses)
		},
	})
}

func (c *V2Context) GetRewardPost(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAggregateMillis),
		utils.NewCounterIncCollect(MetricAggregateCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListTransactionsParams{}
	q, err := ParseGetJSON(r, cfg.RequestGetMaxSize)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	if err := p.ForValues(c.version, q); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	addrsParam := q["addresses"]
	addresses := []string{}
	for _, addrParam := range addrsParam {
		// convert owner address from bech32 to be used internally
		addr, err := address.ParseToID(addrParam)
		if err != nil {
			c.WriteErr(w, 400, err)
			return
		}
		addresses = append(addresses, addr.String())
	}
	// calculate cache key from addresses
	cacheKey := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(addrsParam, ""))))
	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForID("reward", cacheKey),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetReward(ctx, addresses)
		},
	})
}

func (c *V2Context) ListTransactions(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricTransactionsMillis),
		utils.NewCounterIncCollect(MetricTransactionsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListTransactionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)

	if p.ListParams.Offset > DefaultOffsetLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid offset"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_transactions", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListTransactions(ctx, p, c.avaxAssetID)
		},
	})
}

func (c *V2Context) ListTransactionsPost(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricTransactionsMillis),
		utils.NewCounterIncCollect(MetricTransactionsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListTransactionsParams{}
	q, err := ParseGetJSON(r, cfg.RequestGetMaxSize)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	if err := p.ForValues(c.version, q); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)

	if p.ListParams.Offset > DefaultOffsetLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid offset"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_transactions", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListTransactions(ctx, p, c.avaxAssetID)
		},
	})
}

func (c *V2Context) GetTransaction(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricTransactionsMillis),
		utils.NewCounterIncCollect(MetricTransactionsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	id, err := ids.FromString(r.PathParams["id"])
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForID("get_transaction", r.PathParams["id"]),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetTransaction(ctx, id, c.avaxAssetID)
		},
	})
}

func (c *V2Context) ListCTransactions(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricCTransactionsMillis),
		utils.NewCounterIncCollect(MetricCTransactionsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListCTransactionsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	if p.ListParams.Offset > DefaultOffsetLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid offset"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_ctransactions", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListCTransactions(ctx, p)
		},
	})
}

func (c *V2Context) ListCBlocks(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricCTransactionsMillis),
		utils.NewCounterIncCollect(MetricCTransactionsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListCBlocksParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	if p.ListParams.Limit > DefaultLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid block limit"))
		return
	}

	if p.TxLimit > DefaultLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid tx limit"))
		return
	}

	if p.ListParams.Offset > DefaultOffsetLimit {
		c.WriteErr(w, 400, fmt.Errorf("invalid block offset"))
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_cblocks", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListCBlocks(ctx, p)
		},
	})
}

func (c *V2Context) ListAddresses(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAddressesMillis),
		utils.NewCounterIncCollect(MetricAddressesCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListAddressesParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)
	p.ListParams.DisableCounting = true

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_addresses", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListAddresses(ctx, p)
		},
	})
}

func (c *V2Context) GetAddress(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAddressesMillis),
		utils.NewCounterIncCollect(MetricAddressesCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListAddressesParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	id, err := params.AddressFromString(r.PathParams["id"])
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	p.Address = &id
	p.ListParams.DisableCounting = true
	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 1 * time.Second,
		Key: c.cacheKeyForParams("get_address", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetAddress(ctx, p)
		},
	})
}

func (c *V2Context) AddressChains(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAddressChainsMillis),
		utils.NewCounterIncCollect(MetricAddressChainsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.AddressChainsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("address_chains", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.AddressChains(ctx, p)
		},
	})
}

func (c *V2Context) AddressChainsPost(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAddressChainsMillis),
		utils.NewCounterIncCollect(MetricAddressChainsCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.AddressChainsParams{}
	q, err := ParseGetJSON(r, cfg.RequestGetMaxSize)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	if err := p.ForValues(c.version, q); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("address_chains", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.AddressChains(ctx, p)
		},
	})
}

func (c *V2Context) ListOutputs(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListOutputsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	p.ChainIDs = params.ForValueChainID(c.chainID, p.ChainIDs)

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_outputs", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListOutputs(ctx, p)
		},
	})
}

func (c *V2Context) GetOutput(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	id, err := ids.FromString(r.PathParams["id"])
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForID("get_output", r.PathParams["id"]),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetOutput(ctx, id)
		},
	})
}

//
// AVM
//

func (c *V2Context) ListAssets(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAssetMillis),
		utils.NewCounterIncCollect(MetricAssetCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListAssetsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("list_assets", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListAssets(ctx, p, nil)
		},
	})
}

func (c *V2Context) GetAsset(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
		utils.NewCounterObserveMillisCollect(MetricAssetMillis),
		utils.NewCounterIncCollect(MetricAssetCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListAssetsParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	id := r.PathParams["id"]
	p.PathParamID = id

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForParams("get_asset", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetAsset(ctx, p, id)
		},
	})
}

// PVM
func (c *V2Context) ListBlocks(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	p := &params.ListBlocksParams{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		TTL: 5 * time.Second,
		Key: c.cacheKeyForParams("list_blocks", p),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.ListBlocks(ctx, p)
		},
	})
}

func (c *V2Context) GetBlock(w web.ResponseWriter, r *web.Request) {
	collectors := utils.NewCollectors(
		utils.NewCounterObserveMillisCollect(MetricMillis),
		utils.NewCounterIncCollect(MetricCount),
	)
	defer func() {
		_ = collectors.Collect()
	}()

	id, err := ids.FromString(r.PathParams["id"])
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	c.WriteCacheable(w, caching.Cacheable{
		Key: c.cacheKeyForID("get_block", r.PathParams["id"]),
		CacheableFn: func(ctx context.Context) (interface{}, error) {
			return c.avaxReader.GetBlock(ctx, id)
		},
	})
}

func (c *V2Context) ATxData(w web.ResponseWriter, r *web.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	defer cancel()
	p := &params.TxDataParam{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	id := r.PathParams["id"]
	p.ID = id

	b, err := c.avaxReader.ATxDATA(ctx, p)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	WriteJSON(w, b)
}

func (c *V2Context) PTxData(w web.ResponseWriter, r *web.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	defer cancel()
	p := &params.TxDataParam{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	id := r.PathParams["id"]
	p.ID = id

	b, err := c.avaxReader.PTxDATA(ctx, p)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	WriteJSON(w, b)
}

func (c *V2Context) CTxData(w web.ResponseWriter, r *web.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	defer cancel()
	p := &params.TxDataParam{}
	if err := p.ForValues(c.version, r.URL.Query()); err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	id := r.PathParams["id"]
	p.ID = id

	b, err := c.avaxReader.CTxDATA(ctx, p)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}
	WriteJSON(w, b)
}

func (c *V2Context) CacheAddressCounts(w web.ResponseWriter, r *web.Request) {
	res := c.avaxReader.CacheAddressCounts()
	b, err := json.Marshal(res)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	WriteJSON(w, b)
}

func (c *V2Context) CacheTxCounts(w web.ResponseWriter, r *web.Request) {
	res := c.avaxReader.CacheTxCounts()
	b, err := json.Marshal(res)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	WriteJSON(w, b)
}

func (c *V2Context) CacheAssets(w web.ResponseWriter, r *web.Request) {
	res := c.avaxReader.CacheAssets()
	b, err := json.Marshal(res)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	WriteJSON(w, b)
}

func (c *V2Context) CacheAssetAggregates(w web.ResponseWriter, r *web.Request) {
	res := c.avaxReader.CacheAssetAggregates()
	b, err := json.Marshal(res)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	WriteJSON(w, b)
}

func (c *V2Context) CacheAggregates(w web.ResponseWriter, r *web.Request) {
	id := r.PathParams["id"]
	res := c.avaxReader.CacheAggregates(id)
	b, err := json.Marshal(res)
	if err != nil {
		c.WriteErr(w, 400, err)
		return
	}

	WriteJSON(w, b)
}
