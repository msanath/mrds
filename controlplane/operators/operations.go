package operators

import (
	"context"
	"fmt"
	"time"

	"github.com/msanath/gondolf/pkg/ctxslog"
	"github.com/msanath/mrds/controlplane/temporal/workers"
	"github.com/msanath/mrds/controlplane/temporal/workflows"
	"github.com/msanath/mrds/gen/api/mrdspb"
	"go.temporal.io/api/enums/v1"
	temporalclient "go.temporal.io/sdk/client"
)

type operationsOperator struct {
	tc                  temporalclient.Client
	metaInstancesClient mrdspb.MetaInstancesClient
}

func NewOperationsOperator(tc temporalclient.Client, metaInstancesClient mrdspb.MetaInstancesClient) Operator {
	return &operationsOperator{
		tc:                  tc,
		metaInstancesClient: metaInstancesClient,
	}
}

func (d *operationsOperator) RunBlocking(ctx context.Context) error {
	logger := ctxslog.FromContext(ctx)

	ticker, stop := newImmediatelyFiringTicker(10 * time.Second)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled, stopping deployment manager")
			return nil
		case <-ticker:
			// Get all meta instances
			listResp, err := d.metaInstancesClient.List(ctx, &mrdspb.ListMetaInstanceRequest{})
			if err != nil {
				return fmt.Errorf("failed to list meta instances: %w", err)
			}

			for _, metaInstance := range listResp.Records {
				for _, operation := range metaInstance.Operations {
					if operation.Status.State == mrdspb.OperationState_OperationState_PENDING {
						err := d.executeWorkflows(ctx, metaInstance, operation)
						if err != nil {
							return fmt.Errorf("failed to execute workflows: %w", err)
						}
					}
				}
			}
		}
	}
}

func (m *operationsOperator) executeWorkflows(ctx context.Context, metaInstance *mrdspb.MetaInstance, operation *mrdspb.Operation) error {
	log := ctxslog.FromContext(ctx)

	we, err := m.tc.ExecuteWorkflow(ctx,
		temporalclient.StartWorkflowOptions{
			ID:                    fmt.Sprintf("%s-%s", metaInstance.Name, operation.Id),
			TaskQueue:             workers.DeploymentTaskQueue,
			WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		},
		workflows.OperationsWorkflowName,
		&workflows.RunOperationWorkflowParams{
			OperationID:    operation.Id,
			OperationType:  operation.Type,
			MetaInstanceID: metaInstance.Metadata.Id,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}
	log.Info("Started workflow", "workflowID", we.GetID())

	return nil
}
