// Generate an e-invoice from invoice.json and save the PDF and XML. Node.js 18 or newer, no dependencies.
// Usage: INVOIR_API_KEY=inv_... node node/generate.mjs [invoice.json]
import { readFile, writeFile } from 'node:fs/promises';

const key = process.env.INVOIR_API_KEY;
if (!key) throw new Error('Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)');
const headers = { Authorization: `Bearer ${key}` };

const body = await readFile(process.argv[2] ?? 'invoice.json', 'utf8');
const res = await fetch('https://validate.invoir.dev/api/v1/invoices', {
  method: 'POST',
  // Idempotency-Key: a retry with the same key and body returns the first answer instead of a second invoice.
  headers: { ...headers, 'Content-Type': 'application/json', 'Idempotency-Key': JSON.parse(body).number },
  body,
});
const doc = await res.json();
if (!res.ok) {
  for (const e of doc.errors ?? []) console.error(' ', e.field, e.message); // 400: the request
  for (const c of doc.validation?.checks ?? []) { // 422: what the official rules rejected
    for (const f of c.findings) console.error(' ', f.severity, f.rule, f.message, f.location ?? '');
  }
  throw new Error(`${doc.error}: ${doc.message}`);
}

console.log(doc.id, doc.status, doc.totals.due, doc.totals.currency);
for (const file of [doc.pdf, doc.xml].filter(Boolean)) {
  const r = await fetch(file.url, { headers });
  await writeFile(file.filename, Buffer.from(await r.arrayBuffer()));
  console.log('saved', file.filename);
}
