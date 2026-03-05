package browsersigner

// signMessagePageTemplate is the HTML served for personal_sign requests.
const signMessagePageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign Hash with personal_sign</title>
    <style>
        body { font-family: sans-serif; background-color: #1a1a1a; color: #fff; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .container { background: #2d2d2d; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.5); text-align: center; max-width: 600px; width: 100%; }
        button { background-color: #e2761b; color: white; border: none; padding: 15px 30px; font-size: 1.2rem; border-radius: 8px; cursor: pointer; margin-top: 20px; transition: background 0.3s; font-weight: bold; }
        button:hover { background-color: #ff953f; }
        pre { text-align: left; background: #000; padding: 15px; border-radius: 6px; overflow-x: auto; max-height: 300px; font-size: 0.85rem; border: 1px solid #444; word-wrap: break-word; }
        .status { margin-top: 15px; min-height: 24px; color: #aaa; }
        .success { color: #4CAF50; }
        .error { color: #ff5252; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Raw Hash Signing Request</h2>
        <p>Signing the raw 32-byte UserOp hash using <code>personal_sign</code>.</p>

        <div class="status" id="status">Waiting for user action...</div>
        <button id="signBtn">Sign with MetaMask</button>
        <button id="closeBtn" style="display:none" onclick="window.close()">Close Tab</button>

        <h3>Raw Hash (32 Bytes)</h3>
        <pre id="hashDisplay">{{.Payload}}</pre>

    </div>

    <script>
		const urlParams = new URLSearchParams(window.location.search);
		const requestID = urlParams.get('id'); // Get the unique ID

        // Injected Go variable
        const hashToSign = document.getElementById('hashDisplay').textContent.trim();

        const statusDiv = document.getElementById('status');
        const btn = document.getElementById('signBtn');

        btn.addEventListener('click', async () => {
            if (typeof window.ethereum === 'undefined') {
                statusDiv.innerText = "MetaMask is not installed!";
                statusDiv.className = "status error";
                return;
            }

            try {
                statusDiv.innerText = "Requesting account access...";
                // 1. Request accounts
                const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
                const signerAddress = accounts[0];

                statusDiv.innerText = "Requesting signature via personal_sign...";

                // 2. Call personal_sign on the raw hash
                const signature = await window.ethereum.request({
                    method: 'personal_sign',
                    params: [hashToSign, signerAddress], // Message and Address are the required parameters
                });

                statusDiv.innerText = "Signature received! Sending to Go program...";
                statusDiv.className = "status success";

                // 3. Send back to Go server
                const response = await fetch('/submit', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ signature: signature, requestID: requestID })
                });

                if (response.ok) {
                    statusDiv.innerText = "Success! Signature returned.";
                    btn.style.display = 'none';
                    document.getElementById('closeBtn').style.display = 'inline-block';
                } else {
                    statusDiv.innerText = "Server rejected the signature.";
                    statusDiv.className = "status error";
                }
            } catch (err) {
                console.error(err);
                statusDiv.innerText = "Error: " + (err.message || String(err));
                statusDiv.className = "status error";
            }
        });
    </script>
</body>
</html>
`
