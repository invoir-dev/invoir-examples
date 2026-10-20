"""Generate an e-invoice from invoice.json and save the PDF and XML.

Usage: INVOIR_API_KEY=inv_... python3 python/generate.py [invoice.json]
"""
import json
import os
import sys

import requests

key = os.environ.get("INVOIR_API_KEY") or sys.exit("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)")
headers = {"Authorization": f"Bearer {key}"}

with open(sys.argv[1] if len(sys.argv) > 1 else "invoice.json", encoding="utf-8") as f:
    invoice = json.load(f)

# Idempotency-Key: a retry with the same key and body returns the first answer instead of a second invoice.
res = requests.post("https://validate.invoir.dev/api/v1/invoices", json=invoice,
                    headers={**headers, "Idempotency-Key": invoice["number"]}, timeout=120)
doc = res.json()
if not res.ok:
    for e in doc.get("errors", []):  # 400: what is wrong with the request
        print(" ", e["field"], e["message"])
    for check in (doc.get("validation") or {}).get("checks", []):  # 422: what the official rules rejected
        for f in check["findings"]:
            print(" ", f["severity"], f["rule"], f["message"], f.get("location") or "")
    sys.exit(f"{doc['error']}: {doc['message']}")

print(doc["id"], doc["status"], doc["totals"]["due"], doc["totals"]["currency"])
for file in filter(None, [doc["pdf"], doc["xml"]]):
    with open(file["filename"], "wb") as out:
        out.write(requests.get(file["url"], headers=headers, timeout=60).content)
    print("saved", file["filename"])
