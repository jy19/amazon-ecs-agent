// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package v1

import (
	"encoding/json"
	"net/http"

	"github.com/aws/amazon-ecs-agent/agent/config"
	"github.com/aws/amazon-ecs-agent/agent/handlers/utils"
	agentversion "github.com/aws/amazon-ecs-agent/agent/version"
	"github.com/cihub/seelog"
)

// AgentMetadataPath is the Agent metadata path for v1 handler.
const AgentMetadataPath = "/v1/metadata"

// AgentMetadataHandler creates response for 'v1/metadata' API.
func AgentMetadataHandler(containerInstanceArn *string, cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	seelog.Infof("handling v1 metadata for containerInstanceArn %s", *containerInstanceArn)
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &MetadataResponse{
			Cluster:              cfg.Cluster,
			ContainerInstanceArn: containerInstanceArn,
			Version:              agentversion.String(),
		}
		responseJSON, err := json.Marshal(resp)
		if err != nil {
			seelog.Errorf("something went wrong trying to marshal response into JSON, resp: %v", resp)
		}
		utils.WriteJSONToResponse(w, http.StatusOK, responseJSON, utils.RequestTypeAgentMetadata)
	}
}
