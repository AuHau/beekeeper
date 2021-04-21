package config

import (
	"fmt"
	"reflect"
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
	"github.com/ethersphere/beekeeper/pkg/random"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Check struct {
	NewCheck   func() check.Check
	NewOptions func(cfg *Config, checkProfile CheckCfg) (interface{}, error)
}

var Checks = map[string]Check{
	"balances": {
		NewCheck: balances.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				DryRun             *bool   `yaml:"dry-run"`
				FileName           *string `yaml:"file-name"`
				FileSize           *int64  `yaml:"file-size"`
				NodeGroup          *string `yaml:"node-group"`
				Seed               *int64  `yaml:"seed"`
				UploadNodeCount    *int    `yaml:"upload-node-count"`
				WaitBeforeDownload *int    `yaml:"wait-before-download"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := balances.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"chunk-repair": {
		NewCheck: chunkrepair.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				MetricsEnabled         *bool   `yaml:"metrics-enabled"`
				NodeGroup              *string `yaml:"node-group"`
				NumberOfChunksToRepair *int    `yaml:"number-of-chunks-to-repair"`
				Seed                   *int64  `yaml:"seed"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := chunkrepair.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"file-retrieval": {
		NewCheck: fileretrieval.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				FileName        *string `yaml:"file-name"`
				FileSize        *int64  `yaml:"file-size"`
				FilesPerNode    *int    `yaml:"files-per-node"`
				Full            *bool   `yaml:"full"`
				MetricsEnabled  *bool   `yaml:"metrics-enabled"`
				NodeGroup       *string `yaml:"node-group"`
				Seed            *int64  `yaml:"seed"`
				UploadNodeCount *int    `yaml:"upload-node-count"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := fileretrieval.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"full-connectivity": {
		NewCheck: fullconnectivity.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			return nil, nil
		},
	},
	"gc": {
		NewCheck: gc.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				NodeGroup        *string `yaml:"node-group"`
				Seed             *int64  `yaml:"seed"`
				StoreSize        *int    `yaml:"store-size"`
				StoreSizeDivisor *int    `yaml:"store-size-divisor"`
				Wait             *int    `yaml:"wait"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := gc.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"kademlia": {
		NewCheck: kademlia.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				Dynamic *bool `yaml:"dynamic"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := kademlia.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"local-pinning": {
		NewCheck: localpinning.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				Mode             *string `yaml:"mode"`
				NodeGroup        *string `yaml:"node-group"`
				Seed             *int64  `yaml:"seed"`
				StoreSize        *int    `yaml:"store-size"`
				StoreSizeDivisor *int    `yaml:"store-size-divisor"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := localpinning.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"manifest": {
		NewCheck: manifest.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				FilesInCollection *int    `yaml:"files-in-collection"`
				MaxPathnameLength *int32  `yaml:"max-pathname-length"`
				NodeGroup         *string `yaml:"node-group"`
				Seed              *int64  `yaml:"seed"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := manifest.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"peer-count": {
		NewCheck: peercount.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			return nil, nil
		},
	},
	"pingpong": {
		NewCheck: pingpong.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				MetricsEnabled *bool `yaml:"metrics-enabled"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := pingpong.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pss": {
		NewCheck: pss.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				AddressPrefix  *int           `yaml:"address-prefix"`
				MetricsEnabled *bool          `yaml:"metrics-enabled"`
				NodeCount      *int           `yaml:"node-count"`
				NodeGroup      *string        `yaml:"node-group"`
				RequestTimeout *time.Duration `yaml:"request-timeout"`
				Seed           *int64         `yaml:"seed"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := pss.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pullsync": {
		NewCheck: pullsync.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				ChunksPerNode              *int    `yaml:"chunks-per-node"`
				NodeGroup                  *string `yaml:"node-group"`
				ReplicationFactorThreshold *int    `yaml:"replication-factor-threshold"`
				Seed                       *int64  `yaml:"seed"`
				UploadNodeCount            *int    `yaml:"upload-node-count"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := pullsync.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pushsync": {
		NewCheck: pushsync.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
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
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := pushsync.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"retrieval": {
		NewCheck: retrieval.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				ChunksPerNode   *int    `yaml:"chunks-per-node"`
				MetricsEnabled  *bool   `yaml:"metrics-enabled"`
				NodeGroup       *string `yaml:"node-group"`
				Seed            *int64  `yaml:"seed"`
				UploadNodeCount *int    `yaml:"upload-node-count"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := retrieval.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"settlements": {
		NewCheck: settlements.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				DryRun             *bool   `yaml:"dry-run"`
				ExpectSettlements  *bool   `yaml:"expect-settlements"`
				FileName           *string `yaml:"file-name"`
				FileSize           *int64  `yaml:"file-size"`
				NodeGroup          *string `yaml:"node-group"`
				Seed               *int64  `yaml:"seed"`
				Threshold          *int64  `yaml:"threshold"`
				UploadNodeCount    *int    `yaml:"upload-node-count"`
				WaitBeforeDownload *int    `yaml:"wait-before-download"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := settlements.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"soc": {
		NewCheck: soc.NewCheck,
		NewOptions: func(cfg *Config, checkProfile CheckCfg) (interface{}, error) {
			checkOpts := new(struct {
				NodeGroup *string `yaml:"node-group"`
			})
			if err := checkProfile.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkProfile.Name, err)
			}
			opts := soc.NewDefaultOptions()

			if err := applyOptions(cfg, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
}

func applyOptions(cfg *Config, local, opts interface{}) (err error) {
	lv := reflect.ValueOf(local).Elem()
	lt := reflect.TypeOf(local).Elem()
	ov := reflect.Indirect(reflect.ValueOf(opts).Elem())
	ot := reflect.TypeOf(opts).Elem()

	for i := 0; i < lv.NumField(); i++ {
		fieldName := lt.Field(i).Name
		switch fieldName {
		case "MetricsEnabled":
			if lv.Field(i).IsNil() {
				common := cfg.Playbooks[cfg.Execute.Playbook].ChecksCommonOptions.MetricsEnabled
				if common {
					v := reflect.ValueOf(push.New("optionNamePushGateway", cfg.Clusters[cfg.Execute.Cluster].Namespace)) // TODO: resolve optionNamePushGateway
					ov.FieldByName("MetricsPusher").Set(v)
				}
			} else {
				if lv.FieldByName(fieldName).Elem().Bool() {
					v := reflect.ValueOf(push.New("optionNamePushGateway", cfg.Clusters[cfg.Execute.Cluster].Namespace)) // TODO: resolve optionNamePushGateway
					ov.FieldByName("MetricsPusher").Set(v)
				}
			}
		case "Seed":
			if lv.Field(i).IsNil() {
				common := cfg.Playbooks[cfg.Execute.Playbook].ChecksCommonOptions.Seed
				if common >= 0 {
					v := reflect.ValueOf(common)
					ov.FieldByName(fieldName).Set(v)
				} else {
					v := reflect.ValueOf(random.Int64())
					ov.FieldByName(fieldName).Set(v)
				}
			} else {
				fieldType := lt.Field(i).Type
				fieldValue := lv.FieldByName(fieldName).Elem()
				ft, ok := ot.FieldByName(fieldName)
				if ok && fieldType.Elem().AssignableTo(ft.Type) {
					ov.FieldByName(fieldName).Set(fieldValue)
				}
			}
		default:
			if lv.Field(i).IsNil() {
				fmt.Printf("field %s not set\n", fieldName)
			} else {
				fieldType := lt.Field(i).Type
				fieldValue := lv.FieldByName(fieldName).Elem()
				ft, ok := ot.FieldByName(fieldName)
				if ok && fieldType.Elem().AssignableTo(ft.Type) {
					ov.FieldByName(fieldName).Set(fieldValue)
				}
			}
		}
	}

	return
}