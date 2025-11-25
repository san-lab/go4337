package browsersigner

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/san-lab/go4337/userop"
)

// --- MOCK DEPENDENCIES (Replace with your actual structs/functions) ---

// Mock implementation of the assumed HashUserOp function
func HashUserOp(userOp *userop.UserOperation, chainId *big.Int, entrypoint common.Address) [32]byte {
	// In a real scenario, this function would compute the final 32-byte hash
	// using the PackedUserOperation method.
	// We return a hardcoded, valid-looking hash for demonstration.
	hashBytes, _ := hex.DecodeString("1d5420365022e37905f9350280144d4734893116315516087d1587841e248b11")
	var hash [32]byte
	copy(hash[:], hashBytes)
	return hash
}

// --- HTML Template for the browser signing page (Simplified for personal_sign) ---
const signPageTemplate2 = `
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

        <h3>Raw Hash (32 Bytes)</h3>
        <pre id="hashDisplay">{{.Hash}}</pre>
        <div style="font-size: 0.8rem; color: #ffeb3b; margin-top: 10px;">
            ⚠️ MetaMask will automatically prepend a header before signing this hash.
        </div>
    </div>

    <script>
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
                    body: JSON.stringify({ signature: signature })
                });

                if (response.ok) {
                    statusDiv.innerText = "Success! Signature returned. You can close this tab.";
                    btn.style.display = 'none';
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

// ServerConfig holds data for the template
type ServerConfigH struct {
	Hash string
}

// SignHashWithPersonalSign opens a browser to sign the userOp hash using MetaMask's personal_sign.
func SignHashWithPersonalSign(userOp *userop.UserOperation, chainId *big.Int, entrypoint common.Address) ([]byte, error) {

	// 1. Get the hash of the UserOp
	userOpHash := HashUserOp(userOp, chainId, entrypoint)
	userOpHashHex := hexutil.Encode(userOpHash[:])
	fmt.Printf("Generated UserOp Hash: %s\n", userOpHashHex)

	// --- Setup HTTP Server ---
	mux := http.NewServeMux()
	sigChan := make(chan string)

	// Handler for the main signing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("sign").Parse(signPageTemplate2)
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, ServerConfigH{Hash: userOpHashHex})
	})

	// Handler to receive the signature POST
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, _ := io.ReadAll(r.Body)
		var result map[string]string
		if err := json.Unmarshal(body, &result); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Send signature to the channel
		select {
		case sigChan <- result["signature"]:
		default:
		}

		w.WriteHeader(http.StatusOK)
	})

	// Bind to port 0 to find a free random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on local port: %v", err)
	}

	// Get the actual port assigned by the OS
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)

	server := &http.Server{
		Handler: mux,
	}

	// Start Server in Goroutine
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Ensure cleanup happens when function exits
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Printf("Listening on %s\n", url)
	fmt.Println("Opening browser...")

	if err := openBrowser(url); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please open %s manually\n", url)
	}

	fmt.Println("Waiting for signature from Browser...")

	// Wait for signature
	select {
	case signatureHex := <-sigChan:
		// Strip "0x" prefix if present
		if len(signatureHex) >= 2 && signatureHex[:2] == "0x" {
			signatureHex = signatureHex[2:]
		}
		return hex.DecodeString(signatureHex)
	case <-time.After(2 * time.Minute): // Optional timeout
		return nil, fmt.Errorf("timed out waiting for signature")
	}
}
