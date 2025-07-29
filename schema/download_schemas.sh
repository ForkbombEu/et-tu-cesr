#!/usr/bin/env bash
set -euo pipefail

# -------------------------------------------------------------
# 1 – schema file‑names (add/remove freely)
# -------------------------------------------------------------
SCHEMAS=(
  legal-entity-vLEI-credential.json
  legal-entity-engagement-context-role-vLEI-credential.json
  legal-entity-official-organizational-role-vLEI-credential.json
  qualified-vLEI-issuer-vLEI-credential.json
  ecr-authorization-vlei-credential.json
  oor-authorization-vlei-credential.json
  verifiable-ixbrl-report-attestation.json
)

# -------------------------------------------------------------
# 2 – schema versions and where to fetch them
#     • key  = folder name under schema/acdc/
#     • value= *raw* Git URL (no trailing slash)
# -------------------------------------------------------------
declare -A BASE_URLS=(
  # “current”
  [2023]="https://raw.githubusercontent.com/GLEIF-IT/vLEI-schema/main"
  # “old”  (commit 45866b3 on 2022‑06‑24)
  [2022]="https://raw.githubusercontent.com/GLEIF-IT/vLEI-schema/45866b3"
)

# -------------------------------------------------------------
# 3 – download loop
# -------------------------------------------------------------
for ver in "${!BASE_URLS[@]}"; do
  dest="acdc/${ver}"
  mkdir -p "${dest}"

  echo "📂  ${ver}"
  base="${BASE_URLS[$ver]}"

  for f in "${SCHEMAS[@]}"; do
    echo "  ↳ ${f}"
    curl -sSL "${base}/${f}" -o "${dest}/${f}"
  done
done

echo "✅  All schema files downloaded."

