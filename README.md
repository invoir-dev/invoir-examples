# Invoir API examples

Copy-paste examples for the [Invoir](https://invoir.dev) e-invoice API in **curl, Node.js, Python,
PHP and Go**: generate a Factur-X, ZUGFeRD, XRechnung or Peppol BIS invoice from JSON, and validate
any e-invoice against the official rules.

## What Invoir does, in plain words

Across the EU, invoices between businesses are moving from PDFs to **structured e-invoices**: files
that accounting software can read without a person typing anything in. Germany requires businesses to
receive them since 2025 and to send them from 2027 or 2028, France since September 2026 (smaller
businesses from September 2027), and Belgium since January 2026. The formats are all based on the
European standard EN 16931, but each comes with its own rules:

- **Factur-X** / **ZUGFeRD** (the same file under two names): a normal-looking PDF with the invoice
  data embedded as XML. Used in France and Germany.
- **XRechnung**: XML only, the German standard, required by German public authorities.
- **Peppol BIS Billing 3.0**: XML only, sent over the Peppol network, the default in Belgium.

Invoir is an API that does the hard part. You send the invoice as JSON; it computes the totals and
VAT, writes the XML, renders a branded PDF/A-3, **checks the result against the official rule sets**
(KoSIT validator, CEN and Peppol Schematron, Mustang, veraPDF) and only then returns it. It also
validates e-invoices you receive and explains every error with its rule ID. It does not send
invoices: you hand the file to your customer, their platform or your Peppol access point.

Hosted in Germany. The free plan has 25 documents a month; test keys are free and unlimited.

## Quick start

1. Get a free API key at **https://validate.invoir.dev/signup** (email only, no credit card). For
   experiments, create a test key (`inv_test_…`) on the dashboard: test documents are watermarked,
   never counted, and deleted after 24 hours.
2. Put the key in the environment:

   ```bash
   export INVOIR_API_KEY=inv_test_...
   ```

3. From the repository root, generate the sample invoice in [`invoice.json`](invoice.json) and
   validate the PDF you get back:

   ```bash
   node node/generate.mjs          # saves RE-2026-0042.pdf and RE-2026-0042.xml
   node node/validate.mjs RE-2026-0042.pdf
   ```

Every language has the same two scripts, run from the repository root:

| | Generate | Validate | Needs |
| --- | --- | --- | --- |
| curl | `./curl/generate.sh [invoice.json]` | `./curl/validate.sh [file]` | bash, curl, python3 (to read the JSON) |
| Node.js | `node node/generate.mjs [invoice.json]` | `node node/validate.mjs [file]` | Node.js 18 or newer, no packages |
| Python | `python3 python/generate.py [invoice.json]` | `python3 python/validate.py [file]` | `pip install -r python/requirements.txt` |
| PHP | `php php/generate.php [invoice.json]` | `php php/validate.php [file]` | PHP 8 with the curl extension |
| Go | `go run ./go/generate [invoice.json]` | `go run ./go/validate [file]` | Go 1.21 or newer, standard library only |

`generate` posts the JSON to `POST /api/v1/invoices`, prints the document ID, status and amount due,
and saves the PDF and the XML. `validate` posts a file to `POST /api/v1/validate` and prints every
finding with its severity, rule ID, message and location. It accepts a Factur-X or ZUGFeRD PDF, or
XRechnung, Peppol BIS or other EN 16931 XML in UBL or CII.

## The sample invoices

| File | Profile | You get |
| --- | --- | --- |
| [`invoice.json`](invoice.json) | `factur-x` | PDF/A-3 with `factur-x.xml` embedded (EN 16931 profile), plus the XML. Change `profile` to `zugferd` for the same file under the German name. |
| [`invoices/xrechnung-ubl.json`](invoices/xrechnung-ubl.json) | `xrechnung-ubl` | XRechnung 3.0 UBL for a German authority (Leitweg-ID as buyer reference and `0204` address), plus a readable PDF (`"pdf": true`). |
| [`invoices/peppol-bis.json`](invoices/peppol-bis.json) | `peppol-bis` | Peppol BIS Billing 3.0 UBL between two Belgian companies (enterprise numbers, scheme `0208`), plus a readable PDF. |

Companies, VAT IDs, IBANs and the Leitweg-ID in these files are made up (with valid check digits
where the rules check them).

## When something fails

- **400 `invalid_request`**: the JSON can't become an invoice. `errors[]` names every field at once,
  including what the national rules need (Leitweg-ID, electronic addresses, payment details…). The
  scripts print them.
- **422 `invoice_invalid`**: the invoice was generated but the official rules rejected it; the
  scripts print the rule ID and location of each finding from `validation`. Nothing was stored or
  counted.
- **409 `idempotency_conflict`**: the scripts send the invoice number as the `Idempotency-Key`, so a
  retry can never create a second invoice. If you change `invoice.json` and run it again, also change
  `number`.

## Links

- [API documentation](https://invoir.dev/docs) and the [OpenAPI document](https://validate.invoir.dev/openapi.json)
- [Free validator](https://validate.invoir.dev/) in the browser, no account needed
- Guides: [ZUGFeRD](https://invoir.dev/zugferd), [Factur-X](https://invoir.dev/factur-x),
  [XRechnung](https://invoir.dev/xrechnung), [Peppol BIS](https://invoir.dev/peppol-bis),
  [EN 16931](https://invoir.dev/en-16931), and the mandates in
  [Germany](https://invoir.dev/germany-e-invoicing-mandate),
  [France](https://invoir.dev/france-e-invoicing-mandate) and
  [Belgium](https://invoir.dev/belgium-e-invoicing-mandate)
- No code: [n8n, Make and Zapier](https://invoir.dev/integrations)
- [Pricing](https://invoir.dev/pricing) · questions to hello@invoir.dev

## License

MIT. Copy anything you like into your own code.
