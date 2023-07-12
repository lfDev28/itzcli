#!/usr/bin/env bash

# ==============================================================================
# cli-tests.sh
#
# This runs tests against the CLI to verify that the CLI's API does not change
# over time. While there are unit tests that test the code, for the most part,
# these tests are intended to make sure that UNIX command output and exit codes
# are asserted for those people that use `itz` in scripts.
# ==============================================================================
set -ou pipefail

ITZ_WORK_DIR=${ITZ_WORK_DIR:-"$(dirname $0)/../.."}
ITZ_CMD=${ITZ_CMD:-"${ITZ_WORK_DIR}/itzcli"}

log_info() {
  echo $@ >&2
}

fail_msg() {
  echo "$@" 1>&2
  exit 1
}

# Asserts that the provided file exists and is executable
assert_executable() {
  if [ ! -x ${1:-false} ]; then
    fail_msg "    -> Assert failed: expected ${ITZ_CMD} to be executable"
  fi
}

# Asserts both the output and the exit code of the command
assert_code() {
  echo -n "Test output -> ${1}...  "
  cmd_out=$($1 2>&1)
  exit_code=$?
  if [ $exit_code -ne $2 ]; then
    echo "Failed."
    fail_msg "    -> Assert failed: expected ${ITZ_CMD} to have non-error exit code"
  fi
  echo "Passed."
}

# Asserts both the output and the exit code of the command
assert_output_and_code() {
  echo -n "Test output -> ${1}...  "
  cmd_out=$($1 2>&1)
  exit_code=$?
  output_exists=$(echo "$cmd_out" | grep "${2}" | grep -v "grep" | wc -l)
  if [ $output_exists -eq 0 ]; then
    echo "Failed."
    fail_msg "    -> Assert failed: expected ${ITZ_CMD} to have output: \"${2}\""
  fi
  if [ $exit_code -ne $3 ]; then
    echo "Failed."
    fail_msg "    -> Assert failed: expected ${ITZ_CMD} to have non-error exit code"
  fi
  echo "Passed."
}

log_info "Using ${ITZ_CMD} as itz command..."

assert_executable $(command -v ${ITZ_CMD})

# Asserts various commands to make sure that the API (in this case, command structure)
# is stable and doesn't change.
assert_output_and_code "${ITZ_CMD}" "IBM Technology Zone (ITZ) Command Line Interface (CLI)" 0
assert_output_and_code "${ITZ_CMD} execute pipeline" "Error: you must specify a URL for the pipeline to execute" 1
assert_output_and_code "${ITZ_CMD} execute pipeline --pipeline-url moo" "Error: you must specify a URL for the pipeline run to use" 1
assert_output_and_code "${ITZ_CMD} execute pipeline --pipeline-run-url moo" "Error: you must specify a URL for the pipeline to execute" 1

assert_code "${ITZ_CMD} version" 0
# On a fresh system, this should return a non-zero exit code
assert_code "${ITZ_CMD} doctor" 1

exit 0
