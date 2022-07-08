package core

import (
	"github.com/application-research/estuary/pinner/types"
	"github.com/ipfs/go-bitswap"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"time"
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

type FileWithPin struct {
	FileID uint
	PinID  uint

	Cid       string
	Path      string
	Status    types.PinningStatus
	RequestID string
}

type File struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Path      string `gorm:"index"`
	Cid       string
	Mtime     time.Time
}

type Pin struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	File      uint `gorm:"index"`
	Cid       string
	RequestID string `gorm:"index"`
	Status    types.PinningStatus
}
