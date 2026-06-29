#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="${PROTOC_IMAGE:-mtuci-protoc:local}"

mkdir -p "${ROOT}/api/openapi"

PROTO_FILES=()
while IFS= read -r -d '' f; do
  PROTO_FILES+=("proto/${f#"${ROOT}/proto/"}")
done < <(find "${ROOT}/proto" -name '*.proto' -print0)

docker run --rm -v "${ROOT}:/workspace" -w /workspace "${IMAGE}" \
  -I proto \
  --go_out=paths=source_relative:proto \
  --go-grpc_out=paths=source_relative:proto \
  --grpc-gateway_out=paths=source_relative:proto \
  --openapiv2_out=api/openapi \
  "${PROTO_FILES[@]}"

mkdir -p "${ROOT}/internal/server/swaggerdoc"
SWAGGER_JSON="$(find "${ROOT}/api/openapi" -name '*.swagger.json' | head -1)"
if [[ -n "${SWAGGER_JSON}" ]]; then
  cp "${SWAGGER_JSON}" "${ROOT}/internal/server/swaggerdoc/openapi.json"
fi
