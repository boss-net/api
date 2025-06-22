#!/bin/bash
set -x

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
cd "$SCRIPT_DIR/../.."

# ModelRuntime
scripts/pytest/pytest_model_runtime.sh

# Tools
scripts/pytest/pytest_tools.sh

# Workflow
scripts/pytest/pytest_workflow.sh

# Unit tests
scripts/pytest/pytest_unit_tests.sh