#!/bin/bash
set -x

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
cd "$SCRIPT_DIR/../.."

pytest tests/integration_tests/model_runtime/anthropic \
  tests/integration_tests/model_runtime/azure_openai \
  tests/integration_tests/model_runtime/openai tests/integration_tests/model_runtime/chatglm \
  tests/integration_tests/model_runtime/google tests/integration_tests/model_runtime/xinference \
  tests/integration_tests/model_runtime/huggingface_hub/test_llm.py \
  tests/integration_tests/model_runtime/upstage \
  tests/integration_tests/model_runtime/fireworks \
  tests/integration_tests/model_runtime/nomic \
  tests/integration_tests/model_runtime/mixedbread \
  tests/integration_tests/model_runtime/voyage