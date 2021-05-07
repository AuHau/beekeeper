package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func (c *command) initCreateBeeCluster() *cobra.Command {
	const (
		optionNameClusterName = "cluster-name"
		optionNameTimeout     = "timeout"
	)

	cmd := &cobra.Command{
		Use:   "bee-cluster",
		Short: "Create Bee cluster",
		Long:  `Create Bee cluster.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := context.WithTimeout(cmd.Context(), c.globalConfig.GetDuration(optionNameTimeout))
			defer cancel()

			_, err = c.setupCluster(ctx, c.globalConfig.GetString(optionNameClusterName), c.config, true)
			if err != nil {
				return fmt.Errorf("cluster setup: %w", err)
			}

			return
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.globalConfig.BindPFlags(cmd.Flags())
		},
	}

	cmd.Flags().String(optionNameClusterName, "default", "cluster name")
	cmd.Flags().Duration(optionNameTimeout, 15*time.Minute, "timeout")

	c.root.AddCommand(cmd)

	return cmd
}
