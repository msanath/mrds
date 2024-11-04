# Multi-Runtime Deployment System (MRDS)

MRDS is an API-driven deployment system designed to provide a standardized way to interact
with various runtime environments, such as Bare Metal Hosts, Virtual Machines (VMs), Kubernetes Clusters, Azure Virtual Machine Scale Sets, etc.
By abstracting the underlying infrastructure details, MRDS enables developers and operators to
deploy workloads consistently across different runtime systems without needing to tailor
processes to each unique environment.

## Key Features

- **Runtime-Agnostic Interface**: MRDS offers a single, cohesive API that integrates with
diverse runtime environments, enabling smooth deployment operations across platforms.
- **Consistency Across Deployments**: By using MRDS, teams can ensure that deployment
processes are reliable and repeatable, regardless of the runtime specifics.
- **Implementation Flexibility**: Built on a flexible architecture, MRDS can evolve and
scale to accommodate different deployment needs while maintaining a clear, uniform API for users.

## Concepts

To get started with MRDS, it’s important to understand a few core concepts of this system.

### Runtime
In the context of MRDS, a *runtime environment* is the system that manages
workloads on a node. This could be a container orchestration system like Kubernetes or
a virtual machine scale set (VMSS). MRDS is designed to work seamlessly across various
runtime environments, allowing consistent deployment and management processes without
being tied to a specific system.

### Node
A **Node** in MRDS represents an individual compute resource, such as a server or
virtual machine, that is part of a cluster in a runtime environment. Each Node has
specific resources (CPU, memory, storage), a status indicating its current state
(e.g., allocated, unallocated, sanitizing), and may support additional capabilities
or disruptions. MRDS interacts with Nodes through a standardized API to allocate,
monitor, and manage workloads across different runtime systems.


### Deployment Plan

A **Deployment Plan** in MRDS defines what resources and applications should be deployed, along
with specific requirements for the target runtime environment. It is a detailed blueprint
describing the configurations, applications, and compute capabilities needed for deployment.
Each Deployment Plan can be customized to ensure that workloads are deployed only on nodes that
meet specified hardware and software requirements.

An example Deployment Plan may look like this:

```yaml
name: "nginx-deployment-plan"
namespace: "nginx-namespace"
serviceName: "nginx-service"
matchingComputeCapabilities:
    - capabilityType: "GPU"
    comparator: "gte"
    capabilityNames:
        - "nvidia-v100"
        - "nvidia-p100"
    - capabilityType: "CPU"
    comparator: "eq"
    capabilityNames:
        - "intel-xeon"
applications:
    - payloadName: "nginx-payload"
    resources:
        cores: 4
        memory: 64
    ports:
        - protocol: "TCP"
        port: 80
        - protocol: "UDP"
        port: 53
    persistentVolumes:
        - storageClass: "standard"
        capacity: 100
        mountPath: "/data"
        - storageClass: "fast-ssd"
        capacity: 500
        mountPath: "/ssd"
```

### Deployment

A **Deployment** in MRDS is a sub-resource of a Deployment Plan, representing a single execution
of the defined plan. While a Deployment Plan describes the overall configuration and requirements
of an application, a Deployment specifies the actual instances to be launched and managed based on
that plan. Each Deployment can independently control aspects such as the image to deploy and the
number of instances required, allowing for flexibility and scalability.

An example Deployment may look like this:

```yaml
deployment_plan_name: "nginx-deployment-plan"
deployment_id: "deployment-1"
payload_coordinates:
  - payload_name: "nginx-payload"
    coordinates:
      image: "nginx:latest"
instance_count: 1
```
Since multiple Deployments can be associated with a single Deployment Plan, MRDS can support deployment strategies, such as:

- **Blue-Green Deployments**: Run two versions of an application in parallel (one as the live version, and the other as the new version to be switched over).
- **Canary Deployments**: Gradually roll out a new version of the application to a subset of instances before deploying it fully.

### Meta Instance

A **Meta Instance** in MRDS represents a single instance of a Deployment, providing detailed
tracking and management of each deployed unit. Each Meta Instance is uniquely identifiable
and has a 1:1 mapping to an instance of the Deployment.

Meta Instances allow MRDS to manage deployments at a granular level, tracking the precise
state and activity of each deployment instance. By maintaining status, operations, and
runtime instance information, MRDS can effectively control and coordinate multiple instances,
supporting flexible deployment patterns like scaling, rolling updates, and failure handling.

### Runtime Instance

A **Runtime Instance** in MRDS represents the actual instance of an application workload that is
hosted on a specific node within the runtime environment. Each Runtime Instance is associated
with a **Meta Instance** and operates on a particular Node, tracking the execution and state of
the deployed application at the infrastructure level.

Runtime Instances represent the actual compute units for workloads, tracking their status and
activity on specific nodes.

### Operations

**Operations** in MRDS represent specific actions or lifecycle management tasks performed on
Meta Instances or Runtime Instances. Operations allow MRDS to manage instances dynamically,
handling actions such as creation, updates, and terminations based on deployment needs. Each
operation includes information about its type, current status, and progress, making it possible
to track the state of each action performed within the system.

All operations follow a common state transition, ensuring a structured approach to instance
management. The typical process includes:

1. **Creation**: An operation is created in response to a deployment event or change request.
2. **Approval**: Before execution, operations require approval to proceed. Approval can come
   from a manual review or, in many cases, automated systems that evaluate and approve actions
   based on predefined rules and criteria.
3. **Execution**: Once approved, the operation takes action on the instance, performing tasks
   such as deploying, updating, or terminating.
4. **Completion**: After successful execution, the operation is marked as complete. If issues
   arise, it may transition to a `FAILED` state, requiring attention.

Operations allow MRDS to systematically control the full lifecycle of each deployment instance,
adhering to a predictable transition process. With support for both manual and automated
approvals, MRDS can respond flexibly to deployment needs, manage scaling, and adapt instances
based on environmental conditions or predefined rules. This structured approach enables
advanced deployment strategies, such as automated rollouts, scheduled updates, and graceful
shutdowns, enhancing the resilience and responsiveness of deployments.

## Let's deploy

Now that you've been introduced to the core concepts behind MRDS, it's time to put them into
practice. In this section, let's walk through building the system and performing a
deployment. By the end, you'll have a better understanding of MRDS in action, from
defining Deployment Plans to managing individual instances on your chosen runtime environment.


### 1. Install Pre-requisites

For local execution, we are going to use the system to deploy against a locally deployed
kubernetes cluster deployed using [kind](https://kind.sigs.k8s.io/).
- mysql
- [temporal](https://temporal.io/)
- [kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)

The mac commands are
```bash
brew install mysql
brew install temporal
brew install kind
brew install kubectl
```

### Build the system
```bash
make build
```

3 binaries are built and added to the `/bin` directory
- **mrds-apiserver** - The API server for MRDS.
- **mrds-controlplane** - The control plane of MRDS.
- **mrds-ctl** - A CLI for interacting with MRDS.

### Create the database.
```bash
make local-db-reset
```

### Start the API server
Open a new terminal tab to start the API server. Once it's running, you’ll
be able to observe the interactions with the server.
```bash
make run-apiserver
```

### Create a test Kind cluster and add nodes to MRDS.
```bash
make test-kind-prep
```
This creates a local kind cluster with 4 worker nodes. It registers those
nodes to MRDS. You can query them as
```bash
> ./bin/mrds-ctl node list
+--------------+-----------------+-------------+--------------+-----------------------+------------------------+-----------------+------------------+---------------+------------+---------------------+---------+
|  Node Name   |  Update Domain  | Total Cores | Total Memory | System Reserved Cores | System Reserved Memory | Remaining Cores | Remaining Memory | # Disruptions | Cluster ID |        State        | Message |
+--------------+-----------------+-------------+--------------+-----------------------+------------------------+-----------------+------------------+---------------+------------+---------------------+---------+
| kind-worker2 | update-domain-1 |     64      |     512      |           4           |           32           |       60        |       480        |       0       |    kind    | NodeState_ALLOCATED |         |
| kind-worker3 | update-domain-2 |     64      |     512      |           4           |           32           |       60        |       480        |       0       |    kind    | NodeState_ALLOCATED |         |
| kind-worker  | update-domain-1 |     64      |     512      |           4           |           32           |       60        |       480        |       0       |    kind    | NodeState_ALLOCATED |         |
| kind-worker4 | update-domain-2 |     64      |     512      |           4           |           32           |       60        |       480        |       0       |    kind    | NodeState_ALLOCATED |         |
+--------------+-----------------+-------------+--------------+-----------------------+------------------------+-----------------+------------------+---------------+------------+---------------------+---------+
Total: 4
```
(Note: The resources shown here are abstract and not a representation of actual resources.
This setup is intended for demonstration purposes only. In a real runtime environment, these fields will be populated
with actual values.)

### Create a Deployment
The following is a plan to deploy nginx.
```yaml
plans:
  - name: "nginx-deployment-plan"
    namespace: "nginx-namespace"
    serviceName: "nginx-service"
    matchingComputeCapabilities:
      - capabilityType: "GPU"
        comparator: "gte"
        capabilityNames:
          - "nvidia-v100"
          - "nvidia-p100"
      - capabilityType: "CPU"
        comparator: "eq"
        capabilityNames:
          - "intel-xeon"
    applications:
      - payloadName: "nginx-payload"
        resources:
          cores: 4
          memory: 64
        ports:
          - protocol: "TCP"
            port: 80
          - protocol: "UDP"
            port: 53
        persistentVolumes:
          - storageClass: "standard"
            capacity: 100
            mountPath: "/data"
          - storageClass: "fast-ssd"
            capacity: 500
            mountPath: "/ssd"
```
(Note: The capabilities are for illustration only. We are not actually going to deploy on them in this example.)

To create this
```bash
./bin/mrds-ctl deployment create -m test/manifests/deploymentplan.yaml
```

You should see an output like
```bash
Deployment plan created successfullynginx-deployment-plan
+-----------------------+-----------------+---------------+----------------------------+---------+----------------+---------------+
| Deployment Plan Name  |    Namespace    | Service Name  |           State            | Message | # Applications | # Deployments |
+-----------------------+-----------------+---------------+----------------------------+---------+----------------+---------------+
| nginx-deployment-plan | nginx-namespace | nginx-service | DeploymentPlanState_ACTIVE |         |       1        |       0       |
+-----------------------+-----------------+---------------+----------------------------+---------+----------------+---------------+
Total: 1
```

### Add a deployment.
The following is a deployment to deploy 1 instance.
```yaml
deployment_plan_name: "nginx-deployment-plan"
deployment_id: "deployment-1"
payload_coordinates:
  - payload_name: "nginx-payload"
    coordinates:
      image: "nginx:latest"
instance_count: 1
```
Run
```bash
./bin/mrds-ctl deployment add-deployment -m test/manifests/deployment-1.yaml
```
You should see an output as
```bash
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  1
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+-------------------------+---------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |          State          | Message |                        Payload Information                         |
+---------------+----------------+-------------------------+---------+--------------------------------------------------------------------+
| deployment-1  |       1        | DeploymentState_PENDING |         | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+-------------------------+---------+--------------------------------------------------------------------+
Total: 1

Instance Summary
No instances found
```

As mentioned in the output, no instances have been found because there aren't any yet. To create instances, we
need to start the runtime.

### Start temporal server
```bash
make run-temporal
```

### Start controlplane
```bash
make run-controlplane
```

### Observe the changes

Starting the control plane should kickstart a workflow that is going to attempt to deploy
the plan. Open the following temporal link to view the changes - http://localhost:8233/namespaces/mrds/workflows.
Open the workflow with ID - `nginx-deployment-plan-deployment-1` and read through the events.

You'll observe that the workflow created a MetaInstance record and added a CREATE operation on the
instance, and triggered a child workflow. In the child workflow, you'll see that a new `Runtime Instance`
was allocated, and the operation is now waiting approval to proceed.

This can also be observed on the CLI
```bash
> ./bin/mrds-ctl deployment show nginx-deployment-plan
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  2
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+-----------------------------+-----------------------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |            State            |        Message        |                        Payload Information                         |
+---------------+----------------+-----------------------------+-----------------------+--------------------------------------------------------------------+
| deployment-1  |       1        | DeploymentState_IN_PROGRESS | Deployment is running | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+-----------------------------+-----------------------+--------------------------------------------------------------------+
Total: 1

Instance Summary
        Total Instances:          1
        Running Instances:        0
        Pending Instances:        1
        Failed Instances:         0

Opeartions Summary
+------------------------+--------------------+------------+----------+-------------+
|     Operation Type     | # Pending Approval | # Approved | # Failed | # Succeeded |
+------------------------+--------------------+------------+----------+-------------+
|  OperationType_CREATE  |         1          |     0      |    0     |      0      |
|  OperationType_UPDATE  |         0          |     0      |    0     |      0      |
| OperationType_RELOCATE |         0          |     0      |    0     |      0      |
|   OperationType_STOP   |         0          |     0      |    0     |      0      |
| OperationType_RESTART  |         0          |     0      |    0     |      0      |
|  OperationType_DELETE  |         0          |     0      |    0     |      0      |
+------------------------+--------------------+------------+----------+-------------+

Instances
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
|          Meta Instance Name          | Deployment Plan Name  | Deployment ID |            Runtime Instances             |                   Operations                    |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-6P7rl50eSlKsOEbYfWTiwg | nginx-deployment-plan | deployment-1  | Node Name: kind-worker2                  | ID: CREATE-ddd23f06-2050-4515-83e3-a07db62b1b7f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: 8d893aed-a5df-43a9-890c-74677fbe5ff0 | Intent ID: deployment-1                         |
|                                      |                       |               | State: RuntimeState_PENDING              | State: OperationState_PENDING_APPROVAL          |
|                                      |                       |               | Message:                                 | Message: Operation is pending approval          |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
Total: 1
```

As seen above, the node `kind-worker2` has been picked by the scheduler to deploy the workload (This might be different in yourr case).
The State of the RuntimeInstance is `RuntimeStatus_PENDING`, which means, the instance is not  running yet. The last column shows
the operations that are applicable on the meta instance.

### Approve the operation
Execute the following command to select and approve the operations via the CLI.
*Note: These operations are handled through an API, which can also be accessed by automated systems.*
```bash
./bin/mrds-ctl deployment approve-operation nginx-deployment-plan
```

### Observe changes post approval
Keep watching the deployment output on the CLI to see the deployment in action. After a few seconds,
the deployment should be complete, and you should see an output like
```bash
> ./bin/mrds-ctl deployment show nginx-deployment-plan
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  3
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |           State           |              Message              |                        Payload Information                         |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| deployment-1  |       1        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
Total: 1

Instance Summary
        Total Instances:          1
        Running Instances:        1
        Pending Instances:        0
        Failed Instances:         0

Opeartions Summary
+------------------------+--------------------+------------+----------+-------------+
|     Operation Type     | # Pending Approval | # Approved | # Failed | # Succeeded |
+------------------------+--------------------+------------+----------+-------------+
|  OperationType_CREATE  |         0          |     0      |    0     |      1      |
|  OperationType_UPDATE  |         0          |     0      |    0     |      0      |
| OperationType_RELOCATE |         0          |     0      |    0     |      0      |
|   OperationType_STOP   |         0          |     0      |    0     |      0      |
| OperationType_RESTART  |         0          |     0      |    0     |      0      |
|  OperationType_DELETE  |         0          |     0      |    0     |      0      |
+------------------------+--------------------+------------+----------+-------------+

Instances
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
|          Meta Instance Name          | Deployment Plan Name  | Deployment ID |            Runtime Instances             |                   Operations                    |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-6P7rl50eSlKsOEbYfWTiwg | nginx-deployment-plan | deployment-1  | Node Name: kind-worker2                  | ID: CREATE-ddd23f06-2050-4515-83e3-a07db62b1b7f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: 8d893aed-a5df-43a9-890c-74677fbe5ff0 | Intent ID: deployment-1                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
Total: 1
```

You'll see that an instance is RUNNING and the operation was successful. The deployment `deployment-1` has been marked
as complete too.

### Scale up instances
Scaling up instances is a matter of adding a new `Deployment` to the `DeploymentPlan`. Let's increase the instance
count to 3. Run
```bash
 ./bin/mrds-ctl deployment add-deployment -m test/manifests/deployment-2.yaml
```
You should see an output like
```bash
> ./bin/mrds-ctl deployment show nginx-deployment-plan
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  5
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+-----------------------------+-----------------------------------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |            State            |              Message              |                        Payload Information                         |
+---------------+----------------+-----------------------------+-----------------------------------+--------------------------------------------------------------------+
| deployment-1  |       1        |  DeploymentState_COMPLETED  | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
| deployment-2  |       3        | DeploymentState_IN_PROGRESS |       Deployment is running       | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+-----------------------------+-----------------------------------+--------------------------------------------------------------------+
Total: 2

Instance Summary
        Total Instances:          3
        Running Instances:        1
        Pending Instances:        2
        Failed Instances:         0

Opeartions Summary
+------------------------+--------------------+------------+----------+-------------+
|     Operation Type     | # Pending Approval | # Approved | # Failed | # Succeeded |
+------------------------+--------------------+------------+----------+-------------+
|  OperationType_CREATE  |         2          |     0      |    0     |      1      |
|  OperationType_UPDATE  |         1          |     0      |    0     |      0      |
| OperationType_RELOCATE |         0          |     0      |    0     |      0      |
|   OperationType_STOP   |         0          |     0      |    0     |      0      |
| OperationType_RESTART  |         0          |     0      |    0     |      0      |
|  OperationType_DELETE  |         0          |     0      |    0     |      0      |
+------------------------+--------------------+------------+----------+-------------+

Instances
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
|          Meta Instance Name          | Deployment Plan Name  | Deployment ID |            Runtime Instances             |                   Operations                    |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-PXPfTmVLQSmf-Kww_Qbu_w | nginx-deployment-plan | deployment-2  | Node Name: kind-worker3                  | ID: CREATE-f2579172-e4cb-44a5-bd25-389a4522f81f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: c36f2746-a3ae-46bc-8bac-b1c53ad2da87 | Intent ID: deployment-2                         |
|                                      |                       |               | State: RuntimeState_PENDING              | State: OperationState_PENDING_APPROVAL          |
|                                      |                       |               | Message:                                 | Message: Operation is pending approval          |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-S8taiebIQPeJ1sB4y2rS5Q | nginx-deployment-plan | deployment-2  | Node Name: kind-worker                   | ID: CREATE-9f0949d4-a8a2-4315-8f56-c1412c01eec3 |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: dd86c93b-e6fb-43b5-9ceb-90f8171d2421 | Intent ID: deployment-2                         |
|                                      |                       |               | State: RuntimeState_PENDING              | State: OperationState_PENDING_APPROVAL          |
|                                      |                       |               | Message:                                 | Message: Operation is pending approval          |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-6P7rl50eSlKsOEbYfWTiwg | nginx-deployment-plan | deployment-2  | Node Name: kind-worker2                  | ID: CREATE-ddd23f06-2050-4515-83e3-a07db62b1b7f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: 8d893aed-a5df-43a9-890c-74677fbe5ff0 | Intent ID: deployment-1                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
|                                      |                       |               |                                          | ----------------                                |
|                                      |                       |               |                                          | ID: UPDATE-5a4684c0-80b7-461c-baa6-e292b81984c0 |
|                                      |                       |               |                                          | Type: OperationType_UPDATE                      |
|                                      |                       |               |                                          | Intent ID: deployment-2                         |
|                                      |                       |               |                                          | State: OperationState_PENDING_APPROVAL          |
|                                      |                       |               |                                          | Message: Operation is pending approval          |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
Total: 3
```

It shows you that 2 new instances are added, and `CREATE` operations are set on them. The existing one has an operation
`UPDATE` on it, since the instance is already created.

These are again waiting for approval. Approve each one of them and observe the changes.
```bash
./bin/mrds-ctl deployment approve-operation nginx-deployment-plan
```

Once all the approvals are provided, you should see an output like the following after
a few seconds.
```bash
> ./bin/mrds-ctl deployment show nginx-deployment-plan
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  6
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |           State           |              Message              |                        Payload Information                         |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| deployment-1  |       1        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
| deployment-2  |       3        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
Total: 2

Instance Summary
        Total Instances:          3
        Running Instances:        3
        Pending Instances:        0
        Failed Instances:         0

Opeartions Summary
+------------------------+--------------------+------------+----------+-------------+
|     Operation Type     | # Pending Approval | # Approved | # Failed | # Succeeded |
+------------------------+--------------------+------------+----------+-------------+
|  OperationType_CREATE  |         0          |     0      |    0     |      3      |
|  OperationType_UPDATE  |         0          |     0      |    0     |      1      |
| OperationType_RELOCATE |         0          |     0      |    0     |      0      |
|   OperationType_STOP   |         0          |     0      |    0     |      0      |
| OperationType_RESTART  |         0          |     0      |    0     |      0      |
|  OperationType_DELETE  |         0          |     0      |    0     |      0      |
+------------------------+--------------------+------------+----------+-------------+

Instances
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
|          Meta Instance Name          | Deployment Plan Name  | Deployment ID |            Runtime Instances             |                   Operations                    |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-PXPfTmVLQSmf-Kww_Qbu_w | nginx-deployment-plan | deployment-2  | Node Name: kind-worker3                  | ID: CREATE-f2579172-e4cb-44a5-bd25-389a4522f81f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: c36f2746-a3ae-46bc-8bac-b1c53ad2da87 | Intent ID: deployment-2                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-S8taiebIQPeJ1sB4y2rS5Q | nginx-deployment-plan | deployment-2  | Node Name: kind-worker                   | ID: CREATE-9f0949d4-a8a2-4315-8f56-c1412c01eec3 |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: dd86c93b-e6fb-43b5-9ceb-90f8171d2421 | Intent ID: deployment-2                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-6P7rl50eSlKsOEbYfWTiwg | nginx-deployment-plan | deployment-2  | Node Name: kind-worker2                  | ID: CREATE-ddd23f06-2050-4515-83e3-a07db62b1b7f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: 8d893aed-a5df-43a9-890c-74677fbe5ff0 | Intent ID: deployment-1                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
|                                      |                       |               |                                          | ----------------                                |
|                                      |                       |               |                                          | ID: UPDATE-5a4684c0-80b7-461c-baa6-e292b81984c0 |
|                                      |                       |               |                                          | Type: OperationType_UPDATE                      |
|                                      |                       |               |                                          | Intent ID: deployment-2                         |
|                                      |                       |               |                                          | State: OperationState_SUCCEEDED                 |
|                                      |                       |               |                                          | Message:                                        |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
Total: 3
```

Observe that there are 3 corresponding kubernetes pods for the same running on the nodes specified above
```
> kubectl get pods -owide
NAME                                   READY   STATUS    RESTARTS   AGE    IP           NODE           NOMINATED NODE   READINESS GATES
8d893aed-a5df-43a9-890c-74677fbe5ff0   1/1     Running   0          10m    10.244.5.2   kind-worker2   <none>           <none>
c36f2746-a3ae-46bc-8bac-b1c53ad2da87   1/1     Running   0          111s   10.244.3.2   kind-worker3   <none>           <none>
dd86c93b-e6fb-43b5-9ceb-90f8171d2421   1/1     Running   0          110s   10.244.1.2   kind-worker    <none>           <none>
```

### Scale down instance
Scaling down instances is also another deployment. Let's scale down the instances to 1 again.
```bash
./bin/mrds-ctl deployment add-deployment -m test/manifests/deployment-3.yaml
```

Observe the changes again, and approve operations as needed. After a few seconds, you should see an
output like
```bash
> ./bin/mrds-ctl deployment show nginx-deployment-plan
Deployment Plan Info
Deployment Plan Name:     nginx-deployment-plan
Node ID:                  b0ba0333-d386-4625-8743-7b94aa908d32
Version:                  9
Namespace:                nginx-namespace
Service Name:             nginx-service

Status
        Deployment Plan State:    DeploymentPlanState_ACTIVE
        Status Message:           <unset>

Matching Compute Capabilities
+-----------------+---------------+-------------------------+
| Capability Type |  Comparator   |    Capability Names     |
+-----------------+---------------+-------------------------+
|       CPU       | Comparator_IN |       intel-xeon        |
|       GPU       | Comparator_IN | nvidia-v100,nvidia-p100 |
+-----------------+---------------+-------------------------+
Total: 2

Applications
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| Payload Name  | Cores | Memory (GiB) | Ports  |                      Persistent Volumes                      |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
| nginx-payload |   4   |      64      | TCP:80 | Storage Class: standard, Capacity: 100 GB, Mount Path: /data |
|               |       |              | UDP:53 | Storage Class: fast-ssd, Capacity: 500 GB, Mount Path: /ssd  |
+---------------+-------+--------------+--------+--------------------------------------------------------------+
Total: 1

Deployments
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| Deployment ID | Instance Count |           State           |              Message              |                        Payload Information                         |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
| deployment-1  |       1        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
| deployment-2  |       3        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
| deployment-3  |       1        | DeploymentState_COMPLETED | Deployment completed successfully | Payload Name: nginx-payload, Coordinates: {"image":"nginx:latest"} |
+---------------+----------------+---------------------------+-----------------------------------+--------------------------------------------------------------------+
Total: 3

Instance Summary
        Total Instances:          1
        Running Instances:        1
        Pending Instances:        0
        Failed Instances:         0

Opeartions Summary
+------------------------+--------------------+------------+----------+-------------+
|     Operation Type     | # Pending Approval | # Approved | # Failed | # Succeeded |
+------------------------+--------------------+------------+----------+-------------+
|  OperationType_CREATE  |         0          |     0      |    0     |      1      |
|  OperationType_UPDATE  |         0          |     0      |    0     |      2      |
| OperationType_RELOCATE |         0          |     0      |    0     |      0      |
|   OperationType_STOP   |         0          |     0      |    0     |      0      |
| OperationType_RESTART  |         0          |     0      |    0     |      0      |
|  OperationType_DELETE  |         0          |     0      |    0     |      0      |
+------------------------+--------------------+------------+----------+-------------+

Instances
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
|          Meta Instance Name          | Deployment Plan Name  | Deployment ID |            Runtime Instances             |                   Operations                    |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
| nginx-service-6P7rl50eSlKsOEbYfWTiwg | nginx-deployment-plan | deployment-3  | Node Name: kind-worker2                  | ID: CREATE-ddd23f06-2050-4515-83e3-a07db62b1b7f |
|                                      |                       |               | Is Active: true                          | Type: OperationType_CREATE                      |
|                                      |                       |               | ID: 8d893aed-a5df-43a9-890c-74677fbe5ff0 | Intent ID: deployment-1                         |
|                                      |                       |               | State: RuntimeState_RUNNING              | State: OperationState_SUCCEEDED                 |
|                                      |                       |               | Message:                                 | Message:                                        |
|                                      |                       |               |                                          | ----------------                                |
|                                      |                       |               |                                          | ID: UPDATE-5a4684c0-80b7-461c-baa6-e292b81984c0 |
|                                      |                       |               |                                          | Type: OperationType_UPDATE                      |
|                                      |                       |               |                                          | Intent ID: deployment-2                         |
|                                      |                       |               |                                          | State: OperationState_SUCCEEDED                 |
|                                      |                       |               |                                          | Message:                                        |
|                                      |                       |               |                                          | ----------------                                |
|                                      |                       |               |                                          | ID: UPDATE-d4883bc3-afbe-40a8-9566-12525047f308 |
|                                      |                       |               |                                          | Type: OperationType_UPDATE                      |
|                                      |                       |               |                                          | Intent ID: deployment-3                         |
|                                      |                       |               |                                          | State: OperationState_SUCCEEDED                 |
|                                      |                       |               |                                          | Message:                                        |
+--------------------------------------+-----------------------+---------------+------------------------------------------+-------------------------------------------------+
Total: 1
```

And here is the kubectl output
```bash
> kubectl get pods -owide
NAME                                   READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
8d893aed-a5df-43a9-890c-74677fbe5ff0   1/1     Running   0          19m   10.244.5.2   kind-worker2   <none>           <none>
```1

