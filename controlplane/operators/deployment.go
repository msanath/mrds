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

type deploymentOperator struct {
	tc                    temporalclient.Client
	deploymentPlansClient mrdspb.DeploymentPlansClient
}

func NewDeploymentOperator(tc temporalclient.Client, deploymentPlansClient mrdspb.DeploymentPlansClient) Operator {
	return &deploymentOperator{
		tc:                    tc,
		deploymentPlansClient: deploymentPlansClient,
	}
}

func (d *deploymentOperator) RunBlocking(ctx context.Context) error {
	logger := ctxslog.FromContext(ctx)

	ticker, stop := newImmediatelyFiringTicker(10 * time.Second)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled, stopping deployment manager")
			return nil
		case <-ticker:
			// Get all deployment plans
			listResp, err := d.deploymentPlansClient.List(ctx, &mrdspb.ListDeploymentPlanRequest{})
			if err != nil {
				return fmt.Errorf("failed to list deployment plans: %w", err)
			}

			for _, plan := range listResp.Records {
				if plan.Status.State == mrdspb.DeploymentPlanState_DeploymentPlanState_ACTIVE {
					for _, deployment := range plan.Deployments {
						if deployment.Status.State == mrdspb.DeploymentState_DeploymentState_PENDING {
							err := d.executeWorkflows(ctx, plan, deployment)
							if err != nil {
								return fmt.Errorf("failed to execute workflows: %w", err)
							}
						}
					}
				}
			}
		}
	}
}

func (m *deploymentOperator) executeWorkflows(ctx context.Context, deploymentPlan *mrdspb.DeploymentPlanRecord, deployment *mrdspb.Deployment) error {
	log := ctxslog.FromContext(ctx)

	we, err := m.tc.ExecuteWorkflow(ctx,
		temporalclient.StartWorkflowOptions{
			ID:                    fmt.Sprintf("%s-%s", deploymentPlan.Name, deployment.Id),
			TaskQueue:             workers.DeploymentTaskQueue,
			WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		},
		workflows.RunDeploymentWorkflowName,
		&workflows.RunDeploymentWorkflowParams{
			DeploymentPlan: deploymentPlan,
			Deployment:     deployment,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start workflow: %w", err)
	}
	log.Info("Started workflow", "workflowID", we.GetID())

	return nil
}
