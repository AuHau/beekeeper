package cmd

import (
	"fmt"

	"github.com/ethersphere/beekeeper/pkg/config"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"
)

func (c *command) initCheckCmd() (err error) {
	const (
		optionNameClusterName    = "cluster-name"
		optionNameCreateCluster  = "create-cluster"
		optionNameChecks         = "checks"
		optionNameMetricsEnabled = "metrics-enabled"
		optionNameSeed           = "seed"
		// TODO: optionNameStages         = "stages"
		// TODO: optionNameTimeout        = "timeout"
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run tests on a Bee cluster",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// set cluster config
			cfgCluster, ok := c.config.Clusters[c.globalConfig.GetString(optionNameClusterName)]
			if !ok {
				return fmt.Errorf("cluster %s not defined", c.globalConfig.GetString(optionNameClusterName))
			}

			// setup cluster
			cluster, err := c.setupCluster(cmd.Context(), c.globalConfig.GetString(optionNameClusterName), c.config, c.globalConfig.GetBool(optionNameCreateCluster))
			if err != nil {
				return fmt.Errorf("cluster setup: %w", err)
			}

			// set global config
			checkGlobalConfig := config.CheckGlobalConfig{
				MetricsEnabled: c.globalConfig.GetBool(optionNameMetricsEnabled),
				MetricsPusher:  push.New("beekeeper", cfgCluster.GetNamespace()),
				Seed:           c.globalConfig.GetInt64(optionNameSeed),
			}

			// run checks
			for _, checkName := range c.globalConfig.GetStringSlice(optionNameChecks) {
				// get configuration
				checkConfig, ok := c.config.Checks[checkName]
				if !ok {
					return fmt.Errorf("check %s doesn't exist", checkName)
				}

				// choose check type
				check, ok := config.Checks[checkConfig.Type]
				if !ok {
					return fmt.Errorf("check %s not implemented", checkConfig.Type)
				}

				// create check options
				o, err := check.NewOptions(checkGlobalConfig, checkConfig)
				if err != nil {
					return fmt.Errorf("creating check %s options: %w", checkName, err)
				}

				// run check
				if err := check.NewAction().Run(cmd.Context(), cluster, o); err != nil {
					return fmt.Errorf("running check %s: %w", checkName, err)
				}
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.globalConfig.BindPFlags(cmd.Flags())
		},
	}

	cmd.Flags().String(optionNameClusterName, "default", "cluster name")
	cmd.Flags().Bool(optionNameCreateCluster, false, "start cluster")
	cmd.Flags().StringSlice(optionNameChecks, []string{"pingpong"}, "list of checks to execute")
	cmd.Flags().Bool(optionNameMetricsEnabled, false, "enable metrics")
	cmd.Flags().Int64(optionNameSeed, 0, "seed, -1 for random")

	c.root.AddCommand(cmd)

	return nil
}
