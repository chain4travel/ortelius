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

package cfg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/version"
)

const appName = "magellan"

var (
	ErrChainsConfigMustBeStringMap = errors.New("Chain config must a string map")
	ErrChainsConfigIDEmpty         = errors.New("Chain config ID is empty")
	ErrChainsConfigVMEmpty         = errors.New("Chain config vm type is empty")
	ErrChainsConfigIDNotString     = errors.New("Chain config ID is not a string")
	ErrChainsConfigAliasNotString  = errors.New("Chain config alias is not a string")
	ErrChainsConfigVMNotString     = errors.New("Chain config vm type is not a string")
)

type Aggregates struct {
	AggregateMerge    uint64 `json:"AggregateMerge"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	TransactionVolume uint64 `json:"transactionVolume"`
	TransactionCount  uint64 `json:"transactionCount"`
	AddressCount      uint64 `json:"addressCount"`
	OutputCount       uint64 `json:"outputCount"`
	AssetCount        uint64 `json:"assetCount"`
}

type AggregatesMain struct {
	Aggregates Aggregates `json:"aggregates,omitempty"`
	StartTime  string     `json:"startTime"`
	EndTime    string     `json:"endTime"`
}

type AggregatesFees struct {
	AggregateMerge uint64 `json:"AggregateMerge"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	Txfee          uint64 `json:"txfee"`
}

type AggregatesFeesMain struct {
	Aggregates AggregatesFees `json:"aggregates,omitempty"`
	StartTime  string         `json:"startTime"`
	EndTime    string         `json:"endTime"`
}

type Config struct {
	NetworkID               uint32 `json:"networkID"`
	Chains                  `json:"chains"`
	Services                `json:"services"`
	MetricsListenAddr       string `json:"metricsListenAddr"`
	AdminListenAddr         string `json:"adminListenAddr"`
	Features                map[string]struct{}
	CchainID                string `json:"cchainId"`
	CaminoNode              string `json:"caminoNode"`
	NodeInstance            string `json:"nodeInstance"`
	CacheUpdateInterval     uint64 `json:"cacheUpdateInterval"`
	CacheStatisticsInterval uint64 `json:"cacheStatisticsInterval"`
	CacheEmissionsInterval  uint64 `json:"cacheEmissionsInterval"`
	AP5Activation           uint64 `json:"ap5Activation"`
	BanffActivation         uint64 `json:"banffActivation"`
}

type Chain struct {
	ID     string `json:"id"`
	VMType string `json:"vmType"`
}

type Chains map[string]Chain

type Services struct {
	Logging           logging.Config `json:"logging"`
	API               `json:"api"`
	*DB               `json:"db"`
	InmutableInsights EndpointService `json:"inmutableInsights"`
	GeoIP             EndpointService `json:"geoIP"`
}

type EndpointService struct {
	URLEndpoint        string `json:"urlEndpoint"`
	AuthorizationToken string `json:"authorizationToken"`
}

type API struct {
	ListenAddr string `json:"listenAddr"`
}

type DB struct {
	DSN    string `json:"dsn"`
	RODSN  string `json:"rodsn"`
	Driver string `json:"driver"`
}

type Filter struct {
	Min uint32 `json:"min"`
	Max uint32 `json:"max"`
}

// NewFromFile creates a new *Config with the defaults replaced by the config  in
// the file at the given path
func NewFromFile(filePath string) (*Config, error) {
	v, err := newViperFromFile(filePath)
	if err != nil {
		return nil, err
	}

	// Get sub vipers for all objects with parents
	servicesViper := newSubViper(v, keysServices)
	servicesDBViper := newSubViper(servicesViper, keysServicesDB)
	servicesGeoIPViper := newSubViper(servicesViper, keyServicesGeoIP)
	servicesInmutableViper := newSubViper(servicesViper, keyServicesInmutable)

	// Get chains config
	chains, err := newChainsConfig(v)
	if err != nil {
		return nil, err
	}

	// Build logging config
	loggingConf := logging.Config{
		DisplayLevel: logging.Info,
		LogLevel:     logging.Debug,
	}
	loggingConf.Directory = v.GetString(keysLogDirectory)

	dbdsn := servicesDBViper.GetString(keysServicesDBDSN)
	dbrodsn := dbdsn
	if servicesDBViper.Get(keysServicesDBRODSN) != nil {
		dbrodsn = servicesDBViper.GetString(keysServicesDBRODSN)
	}

	urlEndpointGeoIP := servicesGeoIPViper.GetString(keyServicesEndpoint)
	tokenGeoIP := os.Getenv(fmt.Sprintf("%sGeoIP", keyServicesToken))
	urlEndpointInmutable := servicesInmutableViper.GetString(keyServicesEndpoint)
	tokenInmutable := os.Getenv(fmt.Sprintf("%sInmutable", keyServicesToken))

	features := v.GetStringSlice(keysFeatures)
	featuresMap := make(map[string]struct{})
	for _, feature := range features {
		featurec := strings.TrimSpace(strings.ToLower(feature))
		if featurec == "" {
			continue
		}
		featuresMap[featurec] = struct{}{}
	}

	networkID := v.GetUint32(keysNetworkID)
	ap5Activation := version.GetApricotPhase5Time(networkID).Unix()
	banffActivation := version.GetBanffTime(networkID).Unix()

	// Put it all together
	return &Config{
		NetworkID:         networkID,
		Features:          featuresMap,
		Chains:            chains,
		MetricsListenAddr: v.GetString(keysServicesMetricsListenAddr),
		AdminListenAddr:   v.GetString(keysServicesAdminListenAddr),
		Services: Services{
			Logging: loggingConf,
			API: API{
				ListenAddr: v.GetString(keysServicesAPIListenAddr),
			},
			DB: &DB{
				Driver: servicesDBViper.GetString(keysServicesDBDriver),
				DSN:    dbdsn,
				RODSN:  dbrodsn,
			},
			GeoIP: EndpointService{
				URLEndpoint:        urlEndpointGeoIP,
				AuthorizationToken: tokenGeoIP,
			},
			InmutableInsights: EndpointService{
				URLEndpoint:        urlEndpointInmutable,
				AuthorizationToken: tokenInmutable,
			},
		},
		CchainID:                v.GetString(keysStreamProducerCchainID),
		CaminoNode:              v.GetString(keysStreamProducerCaminoNode),
		NodeInstance:            v.GetString(keysStreamProducerNodeInstance),
		CacheUpdateInterval:     uint64(v.GetInt(keysCacheUpdateInterval)),
		CacheStatisticsInterval: uint64(v.GetInt(keysCacheStatisticsInterval)),
		CacheEmissionsInterval:  uint64(v.GetInt(keysCacheEmissionsInterval)),
		AP5Activation:           uint64(ap5Activation),
		BanffActivation:         uint64(banffActivation),
	}, nil
}
