package removepublic

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
	"testing"

	"github.com/google/go-cmp/cmp"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"

	"github.com/googlecloudplatform/security-response-automation/clients/stubs"
	"github.com/googlecloudplatform/security-response-automation/services"
)

func TestCloseCloudSQL(t *testing.T) {
	ctx := context.Background()
	test := []struct {
		name                    string
		instanceDetailsResponse *sqladmin.DatabaseInstance
		expectedRequest         *sqladmin.DatabaseInstance
	}{
		{
			name: "close public ip on sql instance",
			instanceDetailsResponse: &sqladmin.DatabaseInstance{
				Name:    "public-sql-instance",
				Project: "sha-resources-20191002",
				Settings: &sqladmin.Settings{
					IpConfiguration: &sqladmin.IpConfiguration{
						AuthorizedNetworks: []*sqladmin.AclEntry{
							{
								Value: "0.0.0.0/0",
							},
							{
								Value: "199.27.199.0/24",
							},
						},
					},
				},
			},
			expectedRequest: &sqladmin.DatabaseInstance{
				Name:    "public-sql-instance",
				Project: "sha-resources-20191002",
				Settings: &sqladmin.Settings{
					IpConfiguration: &sqladmin.IpConfiguration{
						AuthorizedNetworks: []*sqladmin.AclEntry{
							{
								Value: "199.27.199.0/24",
							},
						},
					},
				},
			},
		},
		{
			name: "tries to close instance already closed",
			instanceDetailsResponse: &sqladmin.DatabaseInstance{
				Name:    "non-public-sql-instance",
				Project: "sha-resources-20191002",
				Settings: &sqladmin.Settings{
					IpConfiguration: &sqladmin.IpConfiguration{
						AuthorizedNetworks: []*sqladmin.AclEntry{
							{
								Value: "199.27.199.0/24",
							},
						},
					},
				},
			},
			expectedRequest: nil,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			svcs, sqlStub := closeSQLSetup()
			sqlStub.InstanceDetailsResponse = tt.instanceDetailsResponse
			values := &Values{
				ProjectID:    "sha-resources-20191002",
				InstanceName: "public-sql-instance",
			}
			if err := Execute(ctx, values, &Services{
				CloudSQL: svcs.CloudSQL,
				Resource: svcs.Resource,
				Logger:   svcs.Logger,
			}); err != nil {
				t.Errorf("%s failed to remove public ip from instance :%q", tt.name, err)
			}

			if diff := cmp.Diff(sqlStub.SavedInstanceUpdated, tt.expectedRequest); diff != "" {
				t.Errorf("%v failed\n exp:%v\n got:%v", tt.name, tt.expectedRequest, sqlStub.SavedInstanceUpdated)
			}
		})
	}
}

func closeSQLSetup() (*services.Global, *stubs.CloudSQL) {
	loggerStub := &stubs.LoggerStub{}
	log := services.NewLogger(loggerStub)
	sqlStub := &stubs.CloudSQL{}
	sql := services.NewCloudSQL(sqlStub)
	storageStub := &stubs.StorageStub{}
	crmStub := &stubs.ResourceManagerStub{}
	res := services.NewResource(crmStub, storageStub)
	return &services.Global{Logger: log, CloudSQL: sql, Resource: res}, sqlStub
}
