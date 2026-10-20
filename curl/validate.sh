#!/usr/bin/env bash
# Validate an e-invoice (Factur-X / ZUGFeRD PDF, XRechnung, Peppol BIS or other EN 16931 XML).
# Usage: INVOIR_API_KEY=inv_... ./curl/validate.sh [file]
set -euo pipefail
: "${INVOIR_API_KEY:?Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)}"
file="${1:-RE-2026-0042.pdf}"

curl -sS https://validate.invoir.dev/api/v1/validate \
  -H "Authorization: Bearer $INVOIR_API_KEY" \
  -F file=@"$file" |
python3 -c '
import json, sys
r = json.load(sys.stdin)
if "error" in r:
    sys.exit(f"{r['"'"'error'"'"']}: {r['"'"'message'"'"']}")
print(r["status"], r["format"], r["errors"], "errors", r["warnings"], "warnings")
for check in r["checks"]:
    for f in check["findings"]:
        print(" ", f["severity"], f["rule"], f["message"], f.get("location") or "")'
