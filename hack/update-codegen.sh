#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

source $(dirname $0)/../vendor/knative.dev/hack/codegen-library.sh
export PATH="$GOBIN:$PATH"

function run_yq() {
	go_run github.com/mikefarah/yq/v4@v4.23.1 "$@"
}

echo "=== Update Codegen for ${MODULE_NAME}"

${CODEGEN_PKG}/generate-groups.sh "client,informer,lister" \
  knative.dev/sample-controller/pkg/avi \
  github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis \
  "ako:v1beta1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

group "Knative Codegen"

# Knative Injection
OUTPUT_PKG="knative.dev/sample-controller/pkg/avi/injection" \
${KNATIVE_CODEGEN_PKG}/hack/generate-knative.sh "injection" \
  github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1beta1 \
  github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis \
  "ako:v1beta1" \
  --go-header-file ${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt

# group "Update deps post-codegen"

# Make sure our dependencies are up-to-date
# ${REPO_ROOT_DIR}/hack/update-deps.sh
