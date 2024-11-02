// Code generated by msanath/gondolf/cligen. DO NOT EDIT.

package types

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/msanath/gondolf/pkg/printer"
)

const (
	ColumnId                 = "id"
	ColumnVersion            = "version"
	ColumnIsDeleted          = "is_deleted"
	ColumnName               = "name"
	ColumnState              = "state"
	ColumnMessage            = "message"
	ColumnDeploymentPlanName = "deployment_plan_name"
	ColumnDeploymentId       = "deployment_id"
)

func GetDisplayMetaInstanceColumnTags() []string {
	return []string{
		ColumnId,
		ColumnVersion,
		ColumnIsDeleted,
		ColumnName,
		ColumnState,
		ColumnMessage,
		ColumnDeploymentPlanName,
		ColumnDeploymentId,
	}
}

func ValidateDisplayMetaInstanceColumnTags(tags []string) error {
	validTags := GetDisplayMetaInstanceColumnTags()
	for _, tag := range tags {
		if !slices.Contains(validTags, tag) {
			return fmt.Errorf("column tag '%s' not found. Valid tags are %v", tag, validTags)
		}
	}
	return nil
}

func (n *DisplayMetaInstance) GetDisplayFieldFromColumnTag(columnTag string) (printer.DisplayField, error) {
	switch columnTag {
	case ColumnId:
		return n.Metadata.GetID(), nil
	case ColumnVersion:
		return n.Metadata.GetVersion(), nil
	case ColumnIsDeleted:
		return n.Metadata.GetIsDeleted(), nil
	case ColumnName:
		return n.GetName(), nil
	case ColumnState:
		return n.Status.GetState(), nil
	case ColumnMessage:
		return n.Status.GetMessage(), nil
	case ColumnDeploymentPlanName:
		return n.GetDeploymentPlanName(), nil
	case ColumnDeploymentId:
		return n.GetDeploymentID(), nil
	}
	return printer.DisplayField{}, fmt.Errorf("column tag '%s' not found. Valid tags are %v", columnTag, GetDisplayMetaInstanceColumnTags())
}

func (n *DisplayMetadata) GetVersion() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Version",
		ColumnTag:   "version",
		Value: func() string {
			str := strconv.Itoa(n.Version)
			return str
		},
	}
}

func (n *DisplayMetadata) GetIsDeleted() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Is Deleted",
		ColumnTag:   "is_deleted",
		Value: func() string {
			str := strconv.FormatBool(n.IsDeleted)
			if str == "true" {
				return printer.RedText(str)
			}
			if str == "false" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DisplayRuntimeInstance) GetIsActive() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Is Active",
		ColumnTag:   "",
		Value: func() string {
			str := strconv.FormatBool(n.IsActive)
			if str == "false" {
				return printer.RedText(str)
			}
			if str == "true" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DisplayMetadata) GetID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "ID",
		ColumnTag:   "id",
		Value: func() string {
			str := n.ID
			return str
		},
	}
}

func (n *DisplayMetaInstance) GetName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Meta Instance Name",
		ColumnTag:   "name",
		Value: func() string {
			str := n.Name
			return str
		},
	}
}

func (n *DisplayMetaInstanceStatus) GetState() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "State",
		ColumnTag:   "state",
		Value: func() string {
			str := n.State
			if str == "MetaInstanceState_TERMINATED" {
				return printer.RedText(str)
			}
			if str == "MetaInstanceState_PENDING_ALLOCATION" {
				return printer.RedText(str)
			}
			if str == "MetaInstanceState_RUNNING" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DisplayMetaInstanceStatus) GetMessage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Status Message",
		ColumnTag:   "message",
		Value: func() string {
			str := n.Message
			return str
		},
	}
}

func (n *DisplayMetaInstance) GetDeploymentPlanName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Deployment Plan Name",
		ColumnTag:   "deployment_plan_name",
		Value: func() string {
			str := n.DeploymentPlanName
			return str
		},
	}
}

func (n *DisplayMetaInstance) GetDeploymentID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Deployment ID",
		ColumnTag:   "deployment_id",
		Value: func() string {
			str := n.DeploymentID
			return str
		},
	}
}

func (n *DisplayRuntimeInstance) GetID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Instance ID",
		ColumnTag:   "",
		Value: func() string {
			str := n.ID
			return str
		},
	}
}

func (n *DisplayRuntimeInstance) GetNodeName() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Node Name",
		ColumnTag:   "",
		Value: func() string {
			str := n.NodeName
			return str
		},
	}
}

func (n *DisplayRuntimeInstanceStatus) GetState() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Runtime State",
		ColumnTag:   "",
		Value: func() string {
			str := n.State
			if str == "RuntimeState_TERMINATED" {
				return printer.RedText(str)
			}
			if str == "RuntimeState_PENDING" {
				return printer.RedText(str)
			}
			if str == "RuntimeState_RUNNING" {
				return printer.GreenText(str)
			}
			return str
		},
	}
}

func (n *DisplayRuntimeInstanceStatus) GetMessage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Status Message",
		ColumnTag:   "",
		Value: func() string {
			str := n.Message
			return str
		},
	}
}

func (n *DisplayOperation) GetID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Operation ID",
		ColumnTag:   "",
		Value: func() string {
			str := n.ID
			return str
		},
	}
}

func (n *DisplayOperation) GetType() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Operation Type",
		ColumnTag:   "",
		Value: func() string {
			str := n.Type
			return str
		},
	}
}

func (n *DisplayOperation) GetIntentID() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Intent ID",
		ColumnTag:   "",
		Value: func() string {
			str := n.IntentID
			return str
		},
	}
}

func (n *DisplayOperationStatus) GetState() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Operation State",
		ColumnTag:   "",
		Value: func() string {
			str := n.State
			if str == "OperationState_FAILED" {
				return printer.RedText(str)
			}
			if str == "OperationState_SUCCEEDED" {
				return printer.GreenText(str)
			}
			if str == "OperationState_PREPARING" {
				return printer.YellowText(str)
			}
			if str == "OperationState_PENDING_APPROVAL" {
				return printer.YellowText(str)
			}
			if str == "OperationState_APPROVED" {
				return printer.YellowText(str)
			}
			return str
		},
	}
}

func (n *DisplayOperationStatus) GetMessage() printer.DisplayField {
	return printer.DisplayField{
		DisplayName: "Status Message",
		ColumnTag:   "",
		Value: func() string {
			str := n.Message
			return str
		},
	}
}
