"""Validate an e-invoice (Factur-X / ZUGFeRD PDF, XRechnung, Peppol BIS or other EN 16931 XML).

Usage: INVOIR_API_KEY=inv_... python3 python/validate.py [file]
"""
import os
import sys

import requests

key = os.environ.get("INVOIR_API_KEY") or sys.exit("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)")
path = sys.argv[1] if len(sys.argv) > 1 else "RE-2026-0042.pdf"

with open(path, "rb") as f:
    res = requests.post(
        "https://validate.invoir.dev/api/v1/validate",
        headers={"Authorization": f"Bearer {key}"},
        files={"file": (os.path.basename(path), f)},
        timeout=60,
    )
report = res.json()
if not res.ok:
    sys.exit(f"{report['error']}: {report['message']}")

print(report["status"], report["format"], report["errors"], "errors", report["warnings"], "warnings")
for check in report["checks"]:
    for finding in check["findings"]:
        print(" ", finding["severity"], finding["rule"], finding["message"], finding.get("location") or "")
