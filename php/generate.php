<?php
// Generate an e-invoice from invoice.json and save the PDF and XML. PHP 8 with the curl extension.
// Usage: INVOIR_API_KEY=inv_... php php/generate.php [invoice.json]
$key = getenv('INVOIR_API_KEY') ?: exit("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)\n");
$invoice = file_get_contents($argv[1] ?? 'invoice.json');
$number = json_decode($invoice, true)['number'];

$ch = curl_init('https://validate.invoir.dev/api/v1/invoices');
curl_setopt_array($ch, [
    CURLOPT_POST => true, CURLOPT_RETURNTRANSFER => true, CURLOPT_TIMEOUT => 120,
    // Idempotency-Key: a retry with the same key and body returns the first answer instead of a second invoice.
    CURLOPT_HTTPHEADER => ["Authorization: Bearer $key", 'Content-Type: application/json', "Idempotency-Key: $number"],
    CURLOPT_POSTFIELDS => $invoice,
]);
$doc = json_decode(curl_exec($ch), true);
$status = curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
if ($status >= 400) {
    foreach ($doc['errors'] ?? [] as $e) echo "  {$e['field']} {$e['message']}\n"; // 400: the request
    foreach ($doc['validation']['checks'] ?? [] as $check) { // 422: what the official rules rejected
        foreach ($check['findings'] as $f) echo "  {$f['severity']} {$f['rule']} {$f['message']} " . ($f['location'] ?? '') . "\n";
    }
    throw new RuntimeException("{$doc['error']}: {$doc['message']}");
}
echo "{$doc['id']} {$doc['status']} {$doc['totals']['due']} {$doc['totals']['currency']}\n";

foreach (array_filter([$doc['pdf'], $doc['xml']]) as $file) {
    $dl = curl_init($file['url']);
    curl_setopt_array($dl, [CURLOPT_RETURNTRANSFER => true, CURLOPT_TIMEOUT => 60, CURLOPT_HTTPHEADER => ["Authorization: Bearer $key"]]);
    file_put_contents($file['filename'], curl_exec($dl));
    echo "saved {$file['filename']}\n";
}
