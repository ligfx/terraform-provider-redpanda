// Copyright 2023 Redpanda Data, Inc.
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ByocCloudProviderDependentValidator is a custom validator to ensure that an attribute is only set when
// cluster_type is BYOC and cloud_provider is a specific value. For example when using this on gcp_project_id
// it will ensure that the HCL fails validation unless cluster_type is set to "byoc" and cloud_provider is
// set to "gcp"
// AttributeName should be the name of the attribute that is being validated
// CloudProvider should be the value of cloud_provider that the attribute is dependent on
type ByocCloudProviderDependentValidator struct {
	AttributeName string
	CloudProvider string
}

// Description provides a description of the validator
func (v ByocCloudProviderDependentValidator) Description(_ context.Context) string {
	return fmt.Sprintf("ensures that %s is only set when cluster_type is byoc and cloud_provider is %s", v.AttributeName, v.CloudProvider)
}

// MarkdownDescription provides a description of the validator in markdown format
func (v ByocCloudProviderDependentValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensures that `%s` is only set when `cluster_type` is `byoc` and `cloud_provider` is `%s`", v.AttributeName, v.CloudProvider)
}

// ValidateString validates a string
func (v ByocCloudProviderDependentValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var clusterType types.String
	var cloudProvider types.String
	if diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("cluster_type"), &clusterType); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("cloud_provider"), &cloudProvider); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// If the object is set and cloud_provider is known but doesn't match, add an error
	if req.ConfigValue.IsNull() {
		return
	}
	if clusterType.IsUnknown() || cloudProvider.IsUnknown() {
		return
	}
	if clusterType.ValueString() != "byoc" || cloudProvider.ValueString() != v.CloudProvider {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Configuration",
			fmt.Sprintf("%s can only be set when cluster_type is byoc and cloud_provider is %s", v.AttributeName, v.CloudProvider),
		)
	}
}
