<?php
// Validate an e-invoice (Factur-X / ZUGFeRD PDF, XRechnung, Peppol BIS or other EN 16931 XML).
// PHP 8 with the curl extension.
// Usage: INVOIR_API_KEY=inv_... php php/validate.php [file]
$key = getenv('INVOIR_API_KEY') ?: exit("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)\n");
$path = $argv[1] ?? 'RE-2026-0042.pdf';

$ch = curl_init('https://validate.invoir.dev/api/v1/validate');
curl_setopt_array($ch, [
    CURLOPT_POST           => true,
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_TIMEOUT        => 60,
    CURLOPT_HTTPHEADER     => ["Authorization: Bearer $key"],
    CURLOPT_POSTFIELDS     => ['file' => new CURLFile($path, '', basename($path))],
]);
$report = json_decode(curl_exec($ch), true);
$status = curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
if ($status >= 400) {
    throw new RuntimeException("{$report['error']}: {$report['message']}");
}

echo "{$report['status']} {$report['format']} {$report['errors']} errors {$report['warnings']} warnings\n";
foreach ($report['checks'] as $check) {
    foreach ($check['findings'] as $f) {
        echo "  {$f['severity']} {$f['rule']} {$f['message']} " . ($f['location'] ?? '') . "\n";
    }
}
