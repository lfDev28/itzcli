#!/usr/bin/env bash
set -euo pipefail

# Print the PATH and the location of the oc command
echo "PATH: $PATH"

ITZ_OC_FLAGS="--insecure-skip-tls-verify=true"
# TODO: add some error handling here
echo "URL: ${ITZ_OC_URL}"

# Try running the oc command with the full path
# Set the OC path to a variable so we can use it later by using which oc
OC_PATH=$(which oc)

$OC_PATH login ${ITZ_OC_URL} -u ${ITZ_OC_USER} -p ${ITZ_OC_PASS}  ${ITZ_OC_FLAGS} 
echo "Applying pipeline..."
$OC_PATH apply -f ${ITZ_PIPELINE}
echo "Applying pipeline run..."
$OC_PATH create -f ${ITZ_PIPELINE_RUN}