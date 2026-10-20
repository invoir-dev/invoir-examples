// Validate an e-invoice (Factur-X / ZUGFeRD PDF, XRechnung, Peppol BIS or other EN 16931 XML).
// Node.js 18 or newer, no dependencies.
// Usage: INVOIR_API_KEY=inv_... node node/validate.mjs [file]
import { readFile } from 'node:fs/promises';
import { basename } from 'node:path';

const key = process.env.INVOIR_API_KEY;
if (!key) throw new Error('Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)');
const path = process.argv[2] ?? 'RE-2026-0042.pdf';

const body = new FormData();
body.append('file', new Blob([await readFile(path)]), basename(path));

const res = await fetch('https://validate.invoir.dev/api/v1/validate', {
  method: 'POST',
  headers: { Authorization: `Bearer ${key}` },
  body,
});
const report = await res.json();
if (!res.ok) throw new Error(`${report.error}: ${report.message}`);

console.log(report.status, report.format, `${report.errors} errors`, `${report.warnings} warnings`);
for (const check of report.checks) {
  for (const f of check.findings) console.log(' ', f.severity, f.rule, f.message, f.location ?? '');
}
