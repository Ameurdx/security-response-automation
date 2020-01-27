// Package clients holds client libraries used by security automation Cloud Functions.
package clients

// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"

	"github.com/googleapis/gax-go/v2"
	commandcenter "github.com/googlecloudplatform/security-response-automation/clients/cscc/apiv1p1alpha1"
	sccpb "github.com/googlecloudplatform/security-response-automation/clients/cscc/v1p1alpha1"
	"google.golang.org/api/option"
)

// SecurityCommandCenter client.
type SecurityCommandCenter struct {
	service *commandcenter.Client
}

// NewSecurityCommandCenter returns and initializes a SecurityCommandCenter client.
func NewSecurityCommandCenter(ctx context.Context, authFile string) (*SecurityCommandCenter, error) {
	cc, err := commandcenter.NewClient(ctx, option.WithCredentialsFile(authFile))
	if err != nil {
		return nil, fmt.Errorf("failed to init scc: %q", err)
	}
	return &SecurityCommandCenter{service: cc}, nil
}

// UpdateSecurityMarks in an Asset or Finding
func (s *SecurityCommandCenter) UpdateSecurityMarks(ctx context.Context, req *sccpb.UpdateSecurityMarksRequest, opts ...gax.CallOption) (*sccpb.SecurityMarks, error) {
	return s.service.UpdateSecurityMarks(ctx, req)
}
