package node

import (
	"context"
	"os"

	"github.com/msanath/mrds/ctl/node/printer"
	"github.com/msanath/mrds/ctl/node/types"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type nodeCreateOptions struct {
	manifestFilePath string

	nodesClient mrdspb.NodesClient
	printer     *printer.Printer
}

type nodeCreateRequest struct {
	Name                    string        `yaml:"name"`
	UpdateDomain            string        `yaml:"updateDomain"`
	TotalResources          resources     `yaml:"totalResources"`
	SystemReservedResources resources     `yaml:"systemReservedResources"`
	CapabilityIDs           []string      `yaml:"capabilityIDs"`
	LocalVolumes            []localVolume `yaml:"localVolumes"`
}

type resources struct {
	Cores  uint32 `yaml:"cores"`
	Memory uint32 `yaml:"memory"`
}

type localVolume struct {
	MountPath       string `yaml:"mountPath"`
	StorageClass    string `yaml:"storageClass"`
	StorageCapacity uint32 `yaml:"storageCapacity"`
}

type nodesList struct {
	Nodes []nodeCreateRequest `yaml:"nodes"`
}

func newNodeCreateCmd() *cobra.Command {
	o := nodeCreateOptions{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new node",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := grpc.NewClient("localhost:12345", grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
			if err != nil {
				return err
			}
			o.nodesClient = mrdspb.NewNodesClient(conn)
			o.printer = printer.NewPrinter()

			return o.run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&o.manifestFilePath, "manifest", "m", "", "Path to the node manifest file")

	return cmd
}

func (o *nodeCreateOptions) run(ctx context.Context) error {
	yamlFile, err := os.Open(o.manifestFilePath)
	if err != nil {
		return err
	}

	req := &nodesList{}
	err = yaml.NewDecoder(yamlFile).Decode(req)
	if err != nil {
		return err
	}

	createdNodes := make([]*mrdspb.Node, 0, len(req.Nodes))
	for _, node := range req.Nodes {
		localVolumes := make([]*mrdspb.NodeLocalVolume, 0, len(node.LocalVolumes))
		for _, lv := range node.LocalVolumes {
			localVolumes = append(localVolumes, &mrdspb.NodeLocalVolume{
				MountPath:       lv.MountPath,
				StorageClass:    lv.StorageClass,
				StorageCapacity: lv.StorageCapacity,
			})
		}

		resp, err := o.nodesClient.Create(ctx, &mrdspb.CreateNodeRequest{
			Name:                    node.Name,
			UpdateDomain:            node.UpdateDomain,
			TotalResources:          &mrdspb.Resources{Cores: node.TotalResources.Cores, Memory: node.TotalResources.Memory},
			SystemReservedResources: &mrdspb.Resources{Cores: node.SystemReservedResources.Cores, Memory: node.SystemReservedResources.Memory},
			CapabilityIds:           node.CapabilityIDs,
			LocalVolumes:            localVolumes,
		})
		if err != nil {
			return err
		}
		createdNodes = append(createdNodes, resp.Record)
	}
	displayNodes := make([]types.DisplayNode, 0, len(createdNodes))
	for _, n := range createdNodes {
		displayNodes = append(displayNodes, convertGRPCNodeToDisplayNode(n))
	}
	o.printer.PrintDisplayNodeList(displayNodes)

	return nil
}
