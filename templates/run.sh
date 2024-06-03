#!/usr/bin/env bash
set -euo pipefail

ITZ_OC_FLAGS="--insecure-skip-tls-verify=true"

OC_PATH=$(command -v oc)

${OC_PATH} login ${ITZ_OC_URL} -u ${ITZ_OC_USER} -p ${ITZ_OC_PASS}  ${ITZ_OC_FLAGS} 
echo "Applying pipeline..."
${OC_PATH} apply -f ${ITZ_PIPELINE}
echo "Applying pipeline run..."
${OC_PATH} create -f ${ITZ_PIPELINE_RUN}