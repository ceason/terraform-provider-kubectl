#!/usr/bin/env bash
[ "$DEBUG" = "1" ] && set -x
set -euo pipefail
err_report() { echo "errexit on line $(caller)" >&2; }
trap err_report ERR

source ${BUILD_WORKING_DIRECTORY:="."}/.rules_terraform/test_vars.sh
kubectl -n$TEST_NAMESPACE get serviceaccount "$TEST_SVCACCT"
kubectl -n$TEST_NAMESPACE get configmap "$TEST_CONFIGMAP"

