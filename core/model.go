package core

import (
	"github.com/ipfs/go-bitswap"
	"github.com/libp2p/go-libp2p-core/host"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
)

type EstuaryConfig struct {
	Token          string `json:"token"`
	Host           string `json:"host"`
	PrimaryShuttle string `json:"primaryShuttle"`
}

type Config struct {
	Estuary EstuaryConfig `json:"estuary"`
}

type EstClient struct {
	Host       string
	Shuttle    string
	Tok        string
	DoProgress bool
	LogTimings bool
}

type PinClient struct {
	host    host.Host
	bitswap *bitswap.Bitswap
	bwc     metrics.Reporter
}
