package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ethersphere/beekeeper/pkg/beekeeper"
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
	"gopkg.in/yaml.v3"
)

type CheckGlobalConfig struct {
	MetricsEnabled bool
	MetricsPusher  *push.Pusher
	Seed           int64
}

type CheckConfig struct {
	Options yaml.Node      `yaml:"options"`
	Timeout *time.Duration `yaml:"timeout"`
	Type    string         `yaml:"type"`
}

type CheckType struct {
	NewAction  func() beekeeper.Action
	NewOptions func(CheckGlobalConfig, CheckConfig) (interface{}, error)
}

var Checks = map[string]CheckType{
	"balances": {
		NewAction: balances.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				DryRun             *bool   `yaml:"dry-run"`
				FileName           *string `yaml:"file-name"`
				FileSize           *int64  `yaml:"file-size"`
				NodeGroup          *string `yaml:"node-group"`
				Seed               *int64  `yaml:"seed"`
				UploadNodeCount    *int    `yaml:"upload-node-count"`
				WaitBeforeDownload *int    `yaml:"wait-before-download"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := balances.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"chunk-repair": {
		NewAction: chunkrepair.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				MetricsEnabled         *bool   `yaml:"metrics-enabled"`
				NodeGroup              *string `yaml:"node-group"`
				NumberOfChunksToRepair *int    `yaml:"number-of-chunks-to-repair"`
				Seed                   *int64  `yaml:"seed"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := chunkrepair.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"file-retrieval": {
		NewAction: fileretrieval.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
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
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := fileretrieval.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"full-connectivity": {
		NewAction: fullconnectivity.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			return nil, nil
		},
	},
	"gc": {
		NewAction: gc.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				NodeGroup        *string `yaml:"node-group"`
				Seed             *int64  `yaml:"seed"`
				StoreSize        *int    `yaml:"store-size"`
				StoreSizeDivisor *int    `yaml:"store-size-divisor"`
				Wait             *int    `yaml:"wait"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := gc.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"kademlia": {
		NewAction: kademlia.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				Dynamic *bool `yaml:"dynamic"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := kademlia.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"local-pinning": {
		NewAction: localpinning.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				Mode             *string `yaml:"mode"`
				NodeGroup        *string `yaml:"node-group"`
				Seed             *int64  `yaml:"seed"`
				StoreSize        *int    `yaml:"store-size"`
				StoreSizeDivisor *int    `yaml:"store-size-divisor"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := localpinning.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"manifest": {
		NewAction: manifest.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				FilesInCollection *int    `yaml:"files-in-collection"`
				MaxPathnameLength *int32  `yaml:"max-pathname-length"`
				NodeGroup         *string `yaml:"node-group"`
				Seed              *int64  `yaml:"seed"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := manifest.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"peer-count": {
		NewAction: peercount.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			return nil, nil
		},
	},
	"pingpong": {
		NewAction: pingpong.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				MetricsEnabled *bool `yaml:"metrics-enabled"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := pingpong.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pss": {
		NewAction: pss.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				AddressPrefix  *int           `yaml:"address-prefix"`
				MetricsEnabled *bool          `yaml:"metrics-enabled"`
				NodeCount      *int           `yaml:"node-count"`
				NodeGroup      *string        `yaml:"node-group"`
				RequestTimeout *time.Duration `yaml:"request-timeout"`
				Seed           *int64         `yaml:"seed"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := pss.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pullsync": {
		NewAction: pullsync.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				ChunksPerNode              *int    `yaml:"chunks-per-node"`
				NodeGroup                  *string `yaml:"node-group"`
				ReplicationFactorThreshold *int    `yaml:"replication-factor-threshold"`
				Seed                       *int64  `yaml:"seed"`
				UploadNodeCount            *int    `yaml:"upload-node-count"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := pullsync.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"pushsync": {
		NewAction: pushsync.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
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
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := pushsync.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"retrieval": {
		NewAction: retrieval.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				ChunksPerNode   *int    `yaml:"chunks-per-node"`
				MetricsEnabled  *bool   `yaml:"metrics-enabled"`
				NodeGroup       *string `yaml:"node-group"`
				Seed            *int64  `yaml:"seed"`
				UploadNodeCount *int    `yaml:"upload-node-count"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := retrieval.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"settlements": {
		NewAction: settlements.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
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
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := settlements.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
	"soc": {
		NewAction: soc.NewCheck,
		NewOptions: func(checkGlobalConfig CheckGlobalConfig, checkConfig CheckConfig) (interface{}, error) {
			checkOpts := new(struct {
				NodeGroup *string `yaml:"node-group"`
			})
			if err := checkConfig.Options.Decode(checkOpts); err != nil {
				return nil, fmt.Errorf("decoding check %s options: %w", checkConfig.Type, err)
			}
			opts := soc.NewDefaultOptions()

			if err := applyCheckConfig(checkGlobalConfig, checkOpts, &opts); err != nil {
				return nil, fmt.Errorf("applying options: %w", err)
			}

			return opts, nil
		},
	},
}

func applyCheckConfig(global CheckGlobalConfig, local, opts interface{}) (err error) {
	lv := reflect.ValueOf(local).Elem()
	lt := reflect.TypeOf(local).Elem()
	ov := reflect.Indirect(reflect.ValueOf(opts).Elem())
	ot := reflect.TypeOf(opts).Elem()

	for i := 0; i < lv.NumField(); i++ {
		fieldName := lt.Field(i).Name
		switch fieldName {
		case "MetricsEnabled":
			// if (set globally) || (set locally)
			if (lv.Field(i).IsNil() && global.MetricsEnabled) || (!lv.Field(i).IsNil() && lv.FieldByName(fieldName).Elem().Bool()) {
				if global.MetricsPusher == nil {
					return fmt.Errorf("metrics pusher is nil (not set)")
				}
				v := reflect.ValueOf(global.MetricsPusher)
				ov.FieldByName("MetricsPusher").Set(v)
			}
		case "Seed":
			if lv.Field(i).IsNil() { // set globally
				if global.Seed >= 0 {
					v := reflect.ValueOf(global.Seed)
					ov.FieldByName(fieldName).Set(v)
				} else {
					v := reflect.ValueOf(random.Int64())
					ov.FieldByName(fieldName).Set(v)
				}
			} else { // set locally
				fieldType := lt.Field(i).Type
				fieldValue := lv.FieldByName(fieldName).Elem()
				ft, ok := ot.FieldByName(fieldName)
				if ok && fieldType.Elem().AssignableTo(ft.Type) {
					ov.FieldByName(fieldName).Set(fieldValue)
				}
			}
		default:
			if lv.Field(i).IsNil() {
				fmt.Printf("field %s not set, using default value\n", fieldName)
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