// Copyright 2023 SAP SE or an SAP affiliate company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spi

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	corev1 "k8s.io/api/core/v1"

	api "github.com/gardener/machine-controller-manager-provider-aws/pkg/aws/apis"
)

// PluginSPIImpl is the real implementation of SPI interface that makes the calls to the AWS SDK.
type PluginSPIImpl struct{}

// NewSession starts a new AWS session
func (ms *PluginSPIImpl) NewSession(secret *corev1.Secret, region string) (*session.Session, error) {
	var config = &aws.Config{
		Region: aws.String(region),
	}

	accessKeyID := extractCredentialsFromData(secret.Data, api.AWSAccessKeyID, api.AWSAlternativeAccessKeyID)
	secretAccessKey := extractCredentialsFromData(secret.Data, api.AWSSecretAccessKey, api.AWSAlternativeSecretAccessKey)

	if accessKeyID != "" && secretAccessKey != "" {
		config = &aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
			}),
		}
	}

	return session.NewSession(config)
}

// NewEC2API Returns a EC2API object
func (ms *PluginSPIImpl) NewEC2API(session *session.Session) ec2iface.EC2API {
	service := ec2.New(session)
	return service
}

// extractCredentialsFromData extracts and trims a value from the given data map. The first key that exists is being
// returned, otherwise, the next key is tried, etc. If no key exists then an empty string is returned.
func extractCredentialsFromData(data map[string][]byte, keys ...string) string {
	for _, key := range keys {
		if val, ok := data[key]; ok {
			return strings.TrimSpace(string(val))
		}
	}
	return ""
}
