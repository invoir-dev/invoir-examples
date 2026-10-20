#!/usr/bin/env bash
# Generate an e-invoice from invoice.json and save the PDF and XML.
# Usage: INVOIR_API_KEY=inv_... ./curl/generate.sh [invoice.json]
set -euo pipefail
: "${INVOIR_API_KEY:?Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)}"
invoice="${1:-invoice.json}"
number=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["number"])' "$invoice")

# Idempotency-Key: a retry with the same key and body returns the first answer instead of a second invoice.
doc=$(curl -sS https://validate.invoir.dev/api/v1/invoices \
  -H "Authorization: Bearer $INVOIR_API_KEY" -H "Content-Type: application/json" \
  -H "Idempotency-Key: $number" --data-binary @"$invoice")

# Print the result, or every problem the API found: 400 lists each field at once, 422 the rules that failed.
echo "$doc" | python3 -c '
import json, sys
d = json.load(sys.stdin)
if "error" in d:
    for e in d.get("errors", []):  # 400: what is wrong with the request
        print(" ", e["field"], e["message"])
    for c in (d.get("validation") or {}).get("checks", []):  # 422: what the official rules rejected
        for f in c["findings"]:
            print(" ", f["severity"], f["rule"], f["message"], f.get("location") or "")
    sys.exit(f"{d['"'"'error'"'"']}: {d['"'"'message'"'"']}")
print(d["id"], d["status"], d["totals"]["due"], d["totals"]["currency"])'

for kind in pdf xml; do
  url=$(echo "$doc" | python3 -c "import json,sys; f=json.load(sys.stdin)['$kind']; print(f['url'] if f else '')")
  name=$(echo "$doc" | python3 -c "import json,sys; f=json.load(sys.stdin)['$kind']; print(f['filename'] if f else '')")
  [ -n "$url" ] || continue
  curl -sS -H "Authorization: Bearer $INVOIR_API_KEY" -o "$name" "$url"
  echo "saved $name"
done
