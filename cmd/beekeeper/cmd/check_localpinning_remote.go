package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ethersphere/beekeeper/pkg/bee"
	"github.com/ethersphere/beekeeper/pkg/check/localpinning"
	"github.com/ethersphere/beekeeper/pkg/random"

	"github.com/spf13/cobra"
)

func (c *command) initCheckLocalPinningRemote() *cobra.Command {
	const (
		optionNameDbCapacity               = "db-capacity"
		optionNameDivisor                  = "capacity-divisor"
		optionNameSeed                     = "seed"
		optionNameStartCluster             = "start-cluster"
		optionNameClusterName              = "cluster-name"
		optionNameBootnodeCount            = "bootnode-count"
		optionNameNodeCount                = "node-count"
		optionNameImage                    = "bee-image"
		optionNamePersistence              = "persistence"
		optionNameStorageClass             = "storage-class"
		optionNameStorageRequest           = "storage-request"
		optionNameFullNode                 = "full-node"
		optionNameAdditionalNodeCount      = "additional-node-count"
		optionNameAdditionalImage          = "additional-bee-image"
		optionNameAdditionalFullNode       = "additional-full-node"
		optionNameAdditionalPersistence    = "additional-persistence"
		optionNameAdditionalStorageClass   = "additional-storage-class"
		optionNameAdditionalStorageRequest = "additional-storage-request"
		optionNameImagePullSecrets         = "image-pull-secrets"
	)

	var (
		imagePullSecrets         []string
		startCluster             bool
		clusterName              string
		bootnodeCount            int
		nodeCount                int
		image                    string
		persistence              bool
		storageClass             string
		storageRequest           string
		fullNode                 bool
		additionalNodeCount      int
		additionalImage          string
		additionalFullNode       bool
		additionalPersistence    bool
		additionalStorageClass   string
		additionalStorageRequest string
	)

	cmdBytes := &cobra.Command{
		Use:   "pin-remote",
		Short: "Checks that a node on the cluster pins remote chunks correctly.",
		Long:  "Checks that a node on the cluster pins remote chunks correctly.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			k8sClient, err := setK8SClient(c.config.GetString(optionNameKubeconfig), c.config.GetBool(optionNameInCluster))
			if err != nil {
				return fmt.Errorf("creating Kubernetes client: %w", err)
			}

			namespace := c.config.GetString(optionNameNamespace)
			cluster := bee.NewCluster(clusterName, bee.ClusterOptions{
				APIDomain:           c.config.GetString(optionNameAPIDomain),
				APIInsecureTLS:      insecureTLSAPI,
				APIScheme:           c.config.GetString(optionNameAPIScheme),
				DebugAPIDomain:      c.config.GetString(optionNameDebugAPIDomain),
				DebugAPIInsecureTLS: insecureTLSDebugAPI,
				DebugAPIScheme:      c.config.GetString(optionNameDebugAPIScheme),
				K8SClient:           k8sClient,
				Namespace:           namespace,
				DisableNamespace:    disableNamespace,
			})

			cicd := newCICDOptions(clefSignerEnable, dbCapacity, paymentEarly, paymentThreshold, paymentTolerance, swapEnable, swapEndpoint, swapFactoryAddress, swapInitialDeposit, nodeSelector, ingressClass)

			if startCluster {
				// bootnodes group
				bgName := "bootnode"
				bCtx, bCancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
				defer bCancel()
				if err := startBootNodeGroup(bCtx, cluster, bootnodeCount, nodeCount, bgName, namespace, image, storageClass, storageRequest, imagePullSecrets, persistence, cicd); err != nil {
					return fmt.Errorf("starting bootnode group %s: %w", bgName, err)
				}

				// node groups
				ngName := "bee"
				nCtx, nCancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
				defer nCancel()
				if err := startNodeGroup(nCtx, cluster, bootnodeCount, nodeCount, ngName, namespace, image, storageClass, storageRequest, imagePullSecrets, persistence, fullNode, cicd); err != nil {
					return fmt.Errorf("starting node group %s: %w", ngName, err)
				}

				if additionalNodeCount > 0 {
					addNgName := "drone"
					addNCtx, addNCancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
					defer addNCancel()
					if err := startNodeGroup(addNCtx, cluster, bootnodeCount, additionalNodeCount, addNgName, namespace, additionalImage, additionalStorageClass, additionalStorageRequest, imagePullSecrets, additionalPersistence, additionalFullNode, cicd); err != nil {
						return fmt.Errorf("starting node group %s: %w", addNgName, err)
					}
				}
			} else {
				// bootnodes group
				if bootnodeCount > 0 {
					bgName := "bootnode"
					if err := addBootNodeGroup(cluster, bootnodeCount, nodeCount, bgName, namespace, image, storageClass, storageRequest, persistence); err != nil {
						return fmt.Errorf("adding bootnode group %s: %w", bgName, err)
					}
				}

				// node groups
				ngName := "bee"
				if err := addNodeGroup(cluster, bootnodeCount, nodeCount, ngName, namespace, image, storageClass, storageRequest, persistence); err != nil {
					return fmt.Errorf("adding node group %s: %w", ngName, err)
				}

				if additionalNodeCount > 0 {
					addNgName := "drone"
					if err := addNodeGroup(cluster, bootnodeCount, additionalNodeCount, addNgName, namespace, additionalImage, additionalStorageClass, additionalStorageRequest, additionalPersistence); err != nil {
						return fmt.Errorf("starting node group %s: %w", addNgName, err)
					}
				}
			}

			var seed int64
			if cmd.Flags().Changed("seed") {
				seed = c.config.GetInt64(optionNameSeed)
			} else {
				seed = random.Int64()
			}

			return localpinning.CheckRemoteChunksFound(cluster, localpinning.Options{
				StoreSize:        c.config.GetInt(optionNameCacheCapacity),
				StoreSizeDivisor: c.config.GetInt(optionNameDivisor),
				Seed:             seed,
				PostageAmount:    c.config.GetInt64(optionNamePostageAmount),
				PostageWait:      c.config.GetDuration(optionNamePostageBatchhWait),
			})
		},
		PreRunE: c.checkPreRunE,
	}

	cmdBytes.Flags().Int(optionNameDivisor, 3, "divide store size by which value when uploading bytes")
	cmdBytes.Flags().Int64P(optionNameSeed, "s", 0, "seed for generating files; if not set, will be random")
	cmdBytes.Flags().BoolVar(&startCluster, optionNameStartCluster, false, "start new cluster")
	cmdBytes.Flags().StringVar(&clusterName, optionNameClusterName, "beekeeper", "cluster name")
	cmdBytes.Flags().IntVarP(&bootnodeCount, optionNameBootnodeCount, "b", 0, "number of bootnodes")
	cmdBytes.Flags().IntVarP(&nodeCount, optionNameNodeCount, "c", 1, "number of nodes")
	cmdBytes.Flags().StringVar(&image, optionNameImage, "ethersphere/bee:latest", "Bee Docker image")
	cmdBytes.PersistentFlags().BoolVar(&persistence, optionNamePersistence, false, "use persistent storage")
	cmdBytes.PersistentFlags().StringVar(&storageClass, optionNameStorageClass, "local-storage", "storage class name")
	cmdBytes.PersistentFlags().StringVar(&storageRequest, optionNameStorageRequest, "34Gi", "storage request")
	cmdBytes.PersistentFlags().BoolVar(&fullNode, optionNameFullNode, true, "start node in full mode")
	cmdBytes.Flags().IntVar(&additionalNodeCount, optionNameAdditionalNodeCount, 0, "number of nodes in additional node group")
	cmdBytes.Flags().StringVar(&additionalImage, optionNameAdditionalImage, "ethersphere/bee:latest", "Bee Docker image in additional node group")
	cmdBytes.PersistentFlags().BoolVar(&additionalFullNode, optionNameAdditionalFullNode, false, "start node in full mode")
	cmdBytes.PersistentFlags().BoolVar(&additionalPersistence, optionNameAdditionalPersistence, false, "use persistent storage")
	cmdBytes.PersistentFlags().StringVar(&additionalStorageClass, optionNameAdditionalStorageClass, "local-storage", "storage class name")
	cmdBytes.PersistentFlags().StringVar(&additionalStorageRequest, optionNameAdditionalStorageRequest, "34Gi", "storage request")
	cmdBytes.Flags().StringArrayVar(&imagePullSecrets, optionNameImagePullSecrets, []string{"regcred"}, "image pull secrets")

	return cmdBytes
}