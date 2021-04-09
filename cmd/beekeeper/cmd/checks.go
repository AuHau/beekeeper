package cmd

import (
	"fmt"
	"time"

	"github.com/ethersphere/beekeeper/pkg/check"
	"github.com/ethersphere/beekeeper/pkg/check/balances"
	"github.com/ethersphere/beekeeper/pkg/check/chunkrepair"
	"github.com/ethersphere/beekeeper/pkg/check/fileretrieval"
	"github.com/ethersphere/beekeeper/pkg/check/fullconnectivity"
	"github.com/ethersphere/beekeeper/pkg/check/gc"
	"github.com/ethersphere/beekeeper/pkg/check/kademlia"
	"github.com/ethersphere/beekeeper/pkg/check/localpinning"
	"github.com/ethersphere/beekeeper/pkg/check/manifest"
	"github.com/ethersphere/beekeeper/pkg/check/peercount"
	"github.com/ethersphere/beekeeper/pkg/check/pingpong"
	"github.com/ethersphere/beekeeper/pkg/check/pss"
	"github.com/ethersphere/beekeeper/pkg/check/pullsync"
	"github.com/ethersphere/beekeeper/pkg/check/pushsync"
	"github.com/ethersphere/beekeeper/pkg/check/retrieval"
	"github.com/ethersphere/beekeeper/pkg/check/settlements"
	"github.com/ethersphere/beekeeper/pkg/check/soc"
	"github.com/ethersphere/beekeeper/pkg/config"
	"github.com/ethersphere/beekeeper/pkg/random"
	"github.com/prometheus/client_golang/prometheus/push"
)

var Checks = map[string]Check{
	"balances": {
		NewCheck: balances.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(balancesOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts balances.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.DryRun != nil {
				opts.DryRun = *o.DryRun
			}
			if o.FileName != nil {
				opts.FileName = *o.FileName
			} else {
				opts.FileName = "balances"
			}
			if o.FileSize != nil {
				opts.FileSize = *o.FileSize
			} else {
				opts.FileSize = 1 * 1024 * 1024 // 1mb
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			if o.WaitBeforeDownload != nil {
				opts.WaitBeforeDownload = *o.WaitBeforeDownload
			} else {
				opts.WaitBeforeDownload = 5 // seconds
			}
			return opts, nil
		},
	},
	"chunk-repair": {
		NewCheck: chunkrepair.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(chunkRepairOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts chunkrepair.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.NumberOfChunksToRepair != nil {
				opts.NumberOfChunksToRepair = *o.NumberOfChunksToRepair
			} else {
				opts.NumberOfChunksToRepair = 1
			}
			return opts, nil
		},
	},
	"file-retrieval": {
		NewCheck: fileretrieval.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(fileRetrievalOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts fileretrieval.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			if o.FileName != nil {
				opts.FileName = *o.FileName
			} else {
				opts.FileName = "file-retrieval"
			}
			if o.FileSize != nil {
				opts.FileSize = *o.FileSize
			} else {
				opts.FileSize = 1 * 1024 * 1024 // 1mb
			}
			if o.FilesPerNode != nil {
				opts.FilesPerNode = *o.FilesPerNode
			} else {
				opts.FilesPerNode = 1
			}
			if o.Full != nil {
				opts.Full = *o.Full
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			return opts, nil
		},
	},
	"full-connectivity": {
		NewCheck: fullconnectivity.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			return nil, nil
		},
	},
	"gc": {
		NewCheck: gc.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(gcOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts gc.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.StoreSize != nil {
				opts.StoreSize = *o.StoreSize
			} else {
				opts.StoreSize = 1000 // DB capacity in chunks
			}
			if o.StoreSizeDivisor != nil {
				opts.StoreSizeDivisor = *o.StoreSizeDivisor
			} else {
				opts.StoreSizeDivisor = 3 // divide store size by which value when uploading bytes
			}
			if o.Wait != nil {
				opts.Wait = *o.Wait
			} else {
				opts.Wait = 5 // wait before check
			}
			return opts, nil
		},
	},
	"kademlia": {
		NewCheck: kademlia.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(kademliaOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts kademlia.Options
			if o.Dynamic != nil {
				opts.Dynamic = *o.Dynamic
			}

			return opts, nil
		},
	},
	"local-pinning": {
		NewCheck: localpinning.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(localpinningOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts localpinning.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.Mode != nil {
				opts.Mode = *o.Mode
			} else {
				opts.Mode = "pin-chunk"
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.StoreSize != nil {
				opts.StoreSize = *o.StoreSize
			} else {
				opts.StoreSize = 1000 // DB capacity in chunks
			}
			if o.StoreSizeDivisor != nil {
				opts.StoreSizeDivisor = *o.StoreSizeDivisor
			} else {
				opts.StoreSizeDivisor = 3 // divide store size by which value when uploading bytes
			}
			return opts, nil
		},
	},
	"manifest": {
		NewCheck: manifest.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(manifestOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			var opts manifest.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.FilesInCollection != nil {
				opts.FilesInCollection = *o.FilesInCollection
			} else {
				opts.FilesInCollection = 10 // number of files to upload in single collection
			}
			if o.MaxPathnameLength != nil {
				opts.MaxPathnameLength = *o.MaxPathnameLength
			} else {
				opts.MaxPathnameLength = 64 // number of files to upload in single collection
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			return opts, nil
		},
	},
	"peer-count": {
		NewCheck: peercount.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			return nil, nil
		},
	},
	"pingpong": {
		NewCheck: pingpong.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(pingOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts pingpong.Options
			// TODO: improve Run["profile"] selection
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			return opts, nil
		},
	},
	"pss": {
		NewCheck: pss.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(pssOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts pss.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			if o.AddressPrefix != nil {
				opts.AddressPrefix = *o.AddressPrefix
			} else {
				opts.AddressPrefix = 1 // public address prefix bytes count
			}
			if o.NodeCount != nil { // TODO: check what this option represent
				opts.NodeCount = *o.NodeCount
			} else {
				opts.NodeCount = 1
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.RequestTimeout != nil {
				opts.RequestTimeout = *o.RequestTimeout
			} else {
				opts.RequestTimeout = 5 * time.Minute
			}
			return opts, nil
		},
	},
	"pullsync": {
		NewCheck: pullsync.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(pullSyncOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts pullsync.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.ChunksPerNode != nil {
				opts.ChunksPerNode = *o.ChunksPerNode
			} else {
				opts.ChunksPerNode = 1 // number of chunks to upload per node
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.ReplicationFactorThreshold != nil {
				opts.ReplicationFactorThreshold = *o.ReplicationFactorThreshold
			} else {
				opts.ReplicationFactorThreshold = 2 // minimal replication factor per chunk
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			return opts, nil
		},
	},
	"pushsync": {
		NewCheck: pushsync.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(pushSyncOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts pushsync.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			if o.ChunksPerNode != nil {
				opts.ChunksPerNode = *o.ChunksPerNode
			} else {
				opts.ChunksPerNode = 1 // number of chunks to upload per node
			}
			if o.FileSize != nil {
				opts.FileSize = *o.FileSize
			} else {
				opts.FileSize = 1 * 1024 * 1024 // 1mb
			}
			if o.FilesPerNode != nil {
				opts.FilesPerNode = *o.FilesPerNode
			} else {
				opts.FilesPerNode = 1
			}
			if o.Mode != nil {
				opts.Mode = *o.Mode
			} else {
				opts.Mode = "default"
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.Retries != nil {
				opts.Retries = *o.Retries
			} else {
				opts.Retries = 5 // number of reties on problems
			}
			if o.RetryDelay != nil {
				opts.RetryDelay = *o.RetryDelay
			} else {
				opts.RetryDelay = 1 * time.Second // retry delay duration
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			return opts, nil
		},
	},
	"retrieval": {
		NewCheck: retrieval.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(retrievalOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts retrieval.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			// TODO: resolve optionNamePushGateway
			// set metrics
			if o.MetricsEnabled == nil && cfg.Run["default"].MetricsEnabled { // enabled globaly
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			} else if o.MetricsEnabled != nil && *o.MetricsEnabled { // enabled localy
				opts.MetricsPusher = push.New("optionNamePushGateway", cfg.Cluster.Namespace)
			}
			if o.ChunksPerNode != nil {
				opts.ChunksPerNode = *o.ChunksPerNode
			} else {
				opts.ChunksPerNode = 1 // number of chunks to upload per node
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			return opts, nil
		},
	},
	"settlements": {
		NewCheck: settlements.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(settlementsOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts settlements.Options
			// TODO: improve Run["profile"] selection
			// set seed
			if o.Seed == nil && cfg.Run["default"].Seed > 0 { // enabled globaly
				opts.Seed = cfg.Run["default"].Seed
			} else if o.Seed != nil && *o.Seed > 0 { // enabled localy
				opts.Seed = *o.Seed
			} else { // randomly generated
				opts.Seed = random.Int64()
			}
			if o.DryRun != nil {
				opts.DryRun = *o.DryRun
			}
			if o.ExpectSettlements != nil {
				opts.ExpectSettlements = *o.ExpectSettlements
			} else {
				opts.ExpectSettlements = true
			}
			if o.FileName != nil {
				opts.FileName = *o.FileName
			} else {
				opts.FileName = "settlements"
			}
			if o.FileSize != nil {
				opts.FileSize = *o.FileSize
			} else {
				opts.FileSize = 1 * 1024 * 1024 // 1mb
			}
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			if o.Threshold != nil {
				opts.Threshold = *o.Threshold
			} else {
				opts.Threshold = 10000000000000 // balances treshold
			}
			if o.UploadNodeCount != nil {
				opts.UploadNodeCount = *o.UploadNodeCount
			} else {
				opts.UploadNodeCount = 1
			}
			if o.WaitBeforeDownload != nil {
				opts.WaitBeforeDownload = *o.WaitBeforeDownload
			} else {
				opts.WaitBeforeDownload = 5 // seconds to wait before downloading a file
			}
			return opts, nil
		},
	},
	"soc": {
		NewCheck: soc.NewCheck,
		NewOptions: func(cfg *config.Config, checkProfile config.Check) (interface{}, error) {
			o := new(socOptions)
			if err := checkProfile.Options.Decode(o); err != nil {
				return nil, fmt.Errorf("decoding check %s optiosns: %w", checkProfile.Name, err)
			}
			var opts soc.Options
			if o.NodeGroup != nil {
				opts.NodeGroup = *o.NodeGroup
			} else {
				opts.NodeGroup = "bee"
			}
			return opts, nil
		},
	},
}

type Check struct {
	NewCheck   func() check.Check
	NewOptions func(cfg *config.Config, checkProfile config.Check) (interface{}, error)
}

type balancesOptions struct {
	DryRun             *bool   `yaml:"dry-run"`
	FileName           *string `yaml:"file-name"`
	FileSize           *int64  `yaml:"file-size"`
	NodeGroup          *string `yaml:"node-group"`
	Seed               *int64  `yaml:"seed"`
	UploadNodeCount    *int    `yaml:"upload-node-count"`
	WaitBeforeDownload *int    `yaml:"wait-before-download"`
}

type chunkRepairOptions struct {
	MetricsEnabled         *bool   `yaml:"metrics-enabled"`
	NodeGroup              *string `yaml:"node-group"`
	NumberOfChunksToRepair *int    `yaml:"number-of-chunks-to-repair"`
	Seed                   *int64  `yaml:"seed"`
}

type fileRetrievalOptions struct {
	FileName        *string `yaml:"file-name"`
	FileSize        *int64  `yaml:"file-size"`
	FilesPerNode    *int    `yaml:"files-per-node"`
	Full            *bool   `yaml:"full"`
	MetricsEnabled  *bool   `yaml:"metrics-enabled"`
	NodeGroup       *string `yaml:"node-group"`
	UploadNodeCount *int    `yaml:"upload-node-count"`
	Seed            *int64  `yaml:"seed"`
}

type gcOptions struct {
	NodeGroup        *string `yaml:"node-group"`
	Seed             *int64  `yaml:"seed"`
	StoreSize        *int    `yaml:"store-size"`
	StoreSizeDivisor *int    `yaml:"store-size-divisor"`
	Wait             *int    `yaml:"wait"`
}

type kademliaOptions struct {
	Dynamic *bool `yaml:"dynamic"`
}

type localpinningOptions struct {
	Mode             *string `yaml:"mode"`
	NodeGroup        *string `yaml:"node-group"`
	Seed             *int64  `yaml:"seed"`
	StoreSize        *int    `yaml:"store-size"`
	StoreSizeDivisor *int    `yaml:"store-size-divisor"`
}

type manifestOptions struct {
	FilesInCollection *int    `yaml:"files-in-collection"`
	MaxPathnameLength *int32  `yaml:"max-pathname-length"`
	NodeGroup         *string `yaml:"node-group"`
	Seed              *int64  `yaml:"seed"`
}

type pingOptions struct {
	MetricsEnabled *bool `yaml:"metrics-enabled"`
}

type pssOptions struct {
	AddressPrefix  *int           `yaml:"address-prefix"`
	MetricsEnabled *bool          `yaml:"metrics-enabled"`
	NodeCount      *int           `yaml:"node-count"`
	NodeGroup      *string        `yaml:"node-group"`
	RequestTimeout *time.Duration `yaml:"request-timeout"`
	Seed           *int64         `yaml:"seed"`
}

type pullSyncOptions struct {
	ChunksPerNode              *int    `yaml:"chunks-per-node"`
	NodeGroup                  *string `yaml:"node-group"`
	ReplicationFactorThreshold *int    `yaml:"replication-factor-threshold"`
	Seed                       *int64  `yaml:"seed"`
	UploadNodeCount            *int    `yaml:"upload-node-count"`
}

type pushSyncOptions struct {
	ChunksPerNode   *int           `yaml:"chunks-per-node"`
	FileSize        *int64         `yaml:"file-size"`
	FilesPerNode    *int           `yaml:"files-per-node"`
	MetricsEnabled  *bool          `yaml:"metrics-enabled"`
	Mode            *string        `yaml:"mode"`
	NodeGroup       *string        `yaml:"node-group"`
	Retries         *int           `yaml:"retries"`
	RetryDelay      *time.Duration `yaml:"retry-delay"`
	Seed            *int64         `yaml:"seed"`
	UploadNodeCount *int           `yaml:"upload-node-count"`
}

type retrievalOptions struct {
	ChunksPerNode   *int    `yaml:"chunks-per-node"`
	MetricsEnabled  *bool   `yaml:"metrics-enabled"`
	NodeGroup       *string `yaml:"node-group"`
	Seed            *int64  `yaml:"seed"`
	UploadNodeCount *int    `yaml:"upload-node-count"`
}

type settlementsOptions struct {
	DryRun             *bool   `yaml:"dry-run"`
	ExpectSettlements  *bool   `yaml:"expect-settlements"`
	FileName           *string `yaml:"file-name"`
	FileSize           *int64  `yaml:"file-size"`
	NodeGroup          *string `yaml:"node-group"`
	Seed               *int64  `yaml:"seed"`
	Threshold          *int64  `yaml:"threshold"`
	UploadNodeCount    *int    `yaml:"upload-node-count"`
	WaitBeforeDownload *int    `yaml:"wait-before-download"`
}

type socOptions struct {
	NodeGroup *string `yaml:"node-group"`
}
