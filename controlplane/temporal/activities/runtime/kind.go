// This file implements a local k8s kind cluster runtime, which is meant to be used for testing purposes.
package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/msanath/mrds/gen/api/mrdspb"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RuntimeActivity interface {
	StartInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error)
	StopInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error)
}

type RuntimeActivityParams struct {
	MetaInstanceID    string
	RuntimeInstanceID string
}

type KindRuntime struct {
	metaInstancesClient  mrdspb.MetaInstancesClient
	deploymentPlanClient mrdspb.DeploymentPlansClient
	nodesClient          mrdspb.NodesClient

	k8sClientSet *kubernetes.Clientset
}

func NewKindRuntime(
	metaInstancesClient mrdspb.MetaInstancesClient,
	deploymentPlanClient mrdspb.DeploymentPlansClient,
	nodesClient mrdspb.NodesClient,
	k8sClientSet *kubernetes.Clientset,
	registry worker.Registry,
) RuntimeActivity {
	a := &KindRuntime{
		metaInstancesClient:  metaInstancesClient,
		deploymentPlanClient: deploymentPlanClient,
		nodesClient:          nodesClient,
		k8sClientSet:         k8sClientSet,
	}

	registry.RegisterActivity(a.StartInstance)
	registry.RegisterActivity(a.StopInstance)

	return a
}

func (k *KindRuntime) StartInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("StartInstance", "request", req)

	runtimeDetails, err := k.getRuntimeDetails(ctx, req.MetaInstanceID, req.RuntimeInstanceID)
	if err != nil {
		return nil, err
	}

	updateResp, err := k.buildAndCreatePod(ctx, runtimeDetails)
	if err != nil {
		return nil, err
	}

	return updateResp, err
}

func (k *KindRuntime) buildAndCreatePod(ctx context.Context, runtimeDetails *runtimeDetails) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("Creating Pod")

	var deployment *mrdspb.Deployment
	for _, d := range runtimeDetails.DeploymentPlan.Deployments {
		if d.Id == runtimeDetails.MetaInstance.DeploymentId {
			deployment = d
			break
		}
	}
	if deployment == nil {
		return nil, fmt.Errorf("deployment with ID %s not found", runtimeDetails.MetaInstance.DeploymentId)
	}

	// build containers for the payloard
	containers := make([]corev1.Container, 0)
	for _, app := range deployment.PayloadCoordinates {
		containers = append(containers, corev1.Container{
			Name:  app.PayloadName,
			Image: app.Coordinates["image"],
		})
	}

	podName := runtimeDetails.RuntimeInstance.Id
	namespace := "default"

	// Define the pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
			NodeSelector: map[string]string{
				"kubernetes.io/hostname": runtimeDetails.Node.Name,
			},
		},
	}

	pod, err := k.k8sClientSet.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			activity.GetLogger(ctx).Info("Pod already exists", "pod", pod)
		} else {
			activity.GetLogger(ctx).Error("Failed to create pod", "error", err)
			return nil, fmt.Errorf("failed to create pod: %v", err)
		}
	} else {
		activity.GetLogger(ctx).Info("Pod created", "pod", pod)
	}

	// Check pod status every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(2 * time.Minute) // Set a timeout period

	for {
		select {
		case <-timeout:
			activity.GetLogger(ctx).Error("Timed out waiting for pod to reach 'Running' state", "pod", pod)
			return nil, fmt.Errorf("timed out waiting for pod %s to reach 'Running' state", podName)
		case <-ticker.C:
			pod, err := k.k8sClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				activity.GetLogger(ctx).Error("Failed to get pod", "error", err)
				continue
			}

			var state mrdspb.RuntimeInstanceState
			switch pod.Status.Phase {
			case corev1.PodPending:
				state = mrdspb.RuntimeInstanceState_RuntimeState_STARTING
			case corev1.PodRunning:
				state = mrdspb.RuntimeInstanceState_RuntimeState_RUNNING
			case corev1.PodSucceeded:
				state = mrdspb.RuntimeInstanceState_RuntimeState_TERMINATED
			case corev1.PodFailed:
				state = mrdspb.RuntimeInstanceState_RuntimeState_FAILED
			case corev1.PodUnknown:
				state = mrdspb.RuntimeInstanceState_RuntimeState_UNKNOWN
			}
			message := pod.Status.Message

			// Update the runtime state to running.
			updateResp, err := k.metaInstancesClient.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
				Metadata:          runtimeDetails.MetaInstance.Metadata,
				RuntimeInstanceId: runtimeDetails.RuntimeInstance.Id,
				Status: &mrdspb.RuntimeInstanceStatus{
					State:   state,
					Message: message,
				},
			})
			if err != nil {
				activity.GetLogger(ctx).Error("Failed to update runtime status", "error", err.Error())
				return nil, err
			}

			// Check if the pod is in the 'Running' phase
			if pod.Status.Phase == corev1.PodRunning {
				activity.GetLogger(ctx).Info("Pod is running", "pod", pod)
				return updateResp, nil
			} else {
				activity.GetLogger(ctx).Info("Pod is not running yet", "pod", pod)
			}
		}
	}
}

func (k *KindRuntime) StopInstance(ctx context.Context, req *RuntimeActivityParams) (*mrdspb.UpdateMetaInstanceResponse, error) {
	activity.GetLogger(ctx).Info("StopInstance", "request", req)

	runtimeDetails, err := k.getRuntimeDetails(ctx, req.MetaInstanceID, req.RuntimeInstanceID)
	if err != nil {
		return nil, err
	}

	// Update the runtime state to stopped.
	updateResp, err := k.metaInstancesClient.UpdateRuntimeStatus(ctx, &mrdspb.UpdateRuntimeStatusRequest{
		Metadata:          runtimeDetails.MetaInstance.Metadata,
		RuntimeInstanceId: req.RuntimeInstanceID,
		Status: &mrdspb.RuntimeInstanceStatus{
			State: mrdspb.RuntimeInstanceState_RuntimeState_TERMINATED,
		},
	})

	return updateResp, err
}

type runtimeDetails struct {
	MetaInstance    *mrdspb.MetaInstance
	RuntimeInstance *mrdspb.RuntimeInstance
	DeploymentPlan  *mrdspb.DeploymentPlanRecord
	Node            *mrdspb.Node
}

func (k *KindRuntime) getRuntimeDetails(ctx context.Context, metaInstanceID string, runtimeInstanceID string) (*runtimeDetails, error) {
	// Get the corresponding meta instance
	metaInstanceGetResp, err := k.metaInstancesClient.GetByID(ctx, &mrdspb.GetMetaInstanceByIDRequest{
		Id: metaInstanceID,
	})
	if err != nil {
		return nil, err
	}

	found := false
	runtimeInstance := &mrdspb.RuntimeInstance{}
	for _, ri := range metaInstanceGetResp.Record.RuntimeInstances {
		if ri.Id == runtimeInstanceID {
			found = true
			runtimeInstance = ri
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("runtime instance %s not found in meta instance %s", metaInstanceID, metaInstanceID)
	}

	// Get the node
	nodeGetResp, err := k.nodesClient.GetByID(ctx, &mrdspb.GetNodeByIDRequest{
		Id: runtimeInstance.NodeId,
	})
	if err != nil {
		return nil, err
	}
	node := nodeGetResp.Record
	activity.GetLogger(ctx).Info("Node", "node", node)

	// Get the corresponding deployment plan
	deploymentPlanGetResp, err := k.deploymentPlanClient.GetByID(ctx, &mrdspb.GetDeploymentPlanByIDRequest{
		Id: metaInstanceGetResp.Record.DeploymentPlanId,
	})
	if err != nil {
		return nil, err
	}
	plan := deploymentPlanGetResp.Record
	activity.GetLogger(ctx).Info("Deployment Plan", "plan", plan)

	return &runtimeDetails{
		MetaInstance:    metaInstanceGetResp.Record,
		RuntimeInstance: runtimeInstance,
		DeploymentPlan:  plan,
		Node:            node,
	}, nil
}
