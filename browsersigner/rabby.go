package browsersigner

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/san-lab/go4337/userop"
)

const EIP7702Prefix = "0x7702000000000000000000000000000000000000"
const EIP7702PrefixLength = 20

// HTML Template for the browser signing page
const signPageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign ERC-4337 v8 UserOp</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #1a1a1a; color: #fff; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; }
        .container { background: #2d2d2d; padding: 2rem; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.5); text-align: center; max-width: 600px; width: 100%; }
        button { background-color: #e2761b; color: white; border: none; padding: 15px 30px; font-size: 1.2rem; border-radius: 8px; cursor: pointer; margin-top: 20px; transition: background 0.3s; font-weight: bold; }
        button:hover { background-color: #ff953f; }
        pre { text-align: left; background: #000; padding: 15px; border-radius: 6px; overflow-x: auto; max-height: 300px; font-size: 0.85rem; border: 1px solid #444; }
        .status { margin-top: 15px; min-height: 24px; color: #aaa; }
        .success { color: #4CAF50; }
        .error { color: #ff5252; }
    </style>
</head>
<body>
    <div class="container">
        <h2>ERC-4337 Signing Request</h2>
        <p>A local Go program requests your signature for the following UserOperation (v8/v0.7).</p>
        
        <div class="status" id="status">Waiting for user action...</div>
        <button id="signBtn">Sign with MetaMask</button>

        <h3>Typed Data Payload</h3>
        <pre id="payloadDisplay"></pre>
    </div>

    <script>
        // Injected Go variable
        const typedDataString = {{.Payload}};
        const typedData = JSON.parse(typedDataString);

        // Display the payload for verification
        document.getElementById('payloadDisplay').textContent = JSON.stringify(typedData, null, 2);

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
                const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
                const signerAddress = accounts[0];

                statusDiv.innerText = "Requesting signature...";
                
                // eth_signTypedData_v4 is the standard for EIP-712
                const signature = await window.ethereum.request({
                    method: 'eth_signTypedData_v4',
                    params: [signerAddress, typedDataString],
                });

                statusDiv.innerText = "Signature received! Sending to Go program...";
                statusDiv.className = "status success";

                // Send back to Go server
                const response = await fetch('/submit', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ signature: signature })
                });

                if (response.ok) {
                    statusDiv.innerText = "Success! You can close this tab.";
                    btn.style.display = 'none';
                } else {
                    statusDiv.innerText = "Server rejected the signature.";
                    statusDiv.className = "status error";
                }
            } catch (err) {
                console.error(err);
                statusDiv.innerText = "Error: " + err.message;
                statusDiv.className = "status error";
            }
        });
    </script>
</body>
</html>
`

// ServerConfig holds data for the template
type ServerConfig struct {
	Payload string
}

func SignEIP712Way(userOp *userop.UserOperation, chainId *big.Int, entrypoint common.Address, eip7702Delegate *common.Address) ([]byte, error) {

	fmt.Println("Generating EIP-712 Payload...")
	payloadJSON, err := PrepareEIP712Payload(userOp, entrypoint, chainId, eip7702Delegate)
	if err != nil {
		return nil, fmt.Errorf("failed to pack UserOp: %v", err)
	}

	// ---------------------------------------------------------
	// FIX 1: Use a new, local ServeMux instead of the global http.HandleFunc
	// This prevents "multiple registrations" panics on subsequent calls.
	// ---------------------------------------------------------
	mux := http.NewServeMux()
	sigChan := make(chan string)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Assuming signPageTemplate and ServerConfig are defined globally or passed in
		tmpl, err := template.New("sign").Parse(signPageTemplate)
		if err != nil {
			http.Error(w, "Template Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, ServerConfig{Payload: payloadJSON})
	})

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

		// Non-blocking send (optional safety, though blocking is usually fine here)
		select {
		case sigChan <- result["signature"]:
		default:
		}

		w.WriteHeader(http.StatusOK)
	})

	// ---------------------------------------------------------
	// FIX 2: Bind to port 0. The OS will assign a free random port.
	// This prevents "address already in use" errors.
	// ---------------------------------------------------------
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

	// Start Server in Goroutine using the existing listener
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

// Helper to open browser cross-platform
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if we are running in WSL (Windows Subsystem for Linux)
		if isWSL() {
			// In WSL, we can call the Windows explorer.exe directly to open the URL
			cmd = "explorer.exe"
		} else {
			// Standard Linux
			cmd = "xdg-open"
		}
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// isWSL checks for the "microsoft" string in /proc/version, which indicates WSL.
func isWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

type TypedData struct {
	Types       map[string][]Type      `json:"types"`
	PrimaryType string                 `json:"primaryType"`
	Domain      Domain                 `json:"domain"`
	Message     map[string]interface{} `json:"message"`
}

type Type struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainId           uint64 `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

// PrepareEIP712Payload formats the UserOp into the JSON string required for browser wallets
func PrepareEIP712Payload(op *userop.UserOperation, entryPoint common.Address, chainID *big.Int, eip7702Delegate *common.Address) (string, error) {

	// 1. PACKING
	accountGasLimits := PackUint128s(op.VerificationGasLimit, op.CallGasLimit)
	gasFees := PackUint128s(op.MaxPriorityFeePerGas, op.MaxFeePerGas)
	paymasterAndData := GetPaymasterAndDataForPacking(op)
	initCode := GetInitCodeForPacking(op, eip7702Delegate)

	// 2. MESSAGE CONSTRUCTION
	message := map[string]interface{}{
		"sender":             op.Sender.Hex(),
		"nonce":              op.Nonce,
		"initCode":           hexutil.Encode(initCode),
		"callData":           hexutil.Encode(op.CallData),
		"accountGasLimits":   hexutil.Encode(accountGasLimits[:]),
		"preVerificationGas": fmt.Sprintf("%d", op.PreVerificationGas),
		"gasFees":            hexutil.Encode(gasFees[:]),
		"paymasterAndData":   hexutil.Encode(paymasterAndData),
	}

	// 3. TYPE DEFINITION
	types := map[string][]Type{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"PackedUserOperation": {
			{Name: "sender", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "initCode", Type: "bytes"},
			{Name: "callData", Type: "bytes"},
			{Name: "accountGasLimits", Type: "bytes32"},
			{Name: "preVerificationGas", Type: "uint256"},
			{Name: "gasFees", Type: "bytes32"},
			{Name: "paymasterAndData", Type: "bytes"},
		},
	}

	// 4. FINAL ASSEMBLY
	typedData := TypedData{
		Types:       types,
		PrimaryType: "PackedUserOperation",
		Domain: Domain{
			Name:              "ERC4337",
			Version:           "1",
			ChainId:           chainID.Uint64(),
			VerifyingContract: entryPoint.Hex(),
		},
		Message: message,
	}

	jsonData, err := json.Marshal(typedData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func GetPaymasterAndDataForPacking(op *userop.UserOperation) []byte {
	if op.Paymaster == nil || *op.Paymaster == (common.Address{}) {
		return []byte{}
	}
	packed := make([]byte, 0, 20+16+16+len(op.PaymasterData))
	packed = append(packed, op.Paymaster.Bytes()...)
	verifGasBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(verifGasBytes[8:], op.PaymasterVerificationGasLimit)
	packed = append(packed, verifGasBytes...)
	postOpBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(postOpBytes[8:], op.PaymasterPostOpGasLimit)
	packed = append(packed, postOpBytes...)
	packed = append(packed, op.PaymasterData...)
	return packed
}

func GetInitCodeForPacking(op *userop.UserOperation, eip7702Delegate *common.Address) []byte {
	var baseInitCode []byte
	if op.Factory != nil && *op.Factory != (common.Address{}) {
		baseInitCode = append(op.Factory.Bytes(), op.FactoryData...)
	} else if len(op.FactoryData) > 0 {
		// This path covers cases where op.InitCode is set directly (non-factory)
		baseInitCode = op.FactoryData
	}

	// *** MODIFICATION START ***
	// Apply EIP-7702 Hash Tweak: Check for the special prefix
	if len(baseInitCode) >= EIP7702PrefixLength && hexutil.Encode(baseInitCode[:EIP7702PrefixLength]) == EIP7702Prefix {
		// Copy the base initCode to modify it for the hash calculation
		modifiedInitCode := make([]byte, len(baseInitCode))
		copy(modifiedInitCode, baseInitCode)

		// Replace the first 20 bytes (the prefix) with the EIP-7702 Delegate Address
		copy(modifiedInitCode[:EIP7702PrefixLength], eip7702Delegate.Bytes())
		return modifiedInitCode
	}
	return baseInitCode
}

func PackUint128s(high, low uint64) [32]byte {
	var packed [32]byte
	binary.BigEndian.PutUint64(packed[8:16], high)
	binary.BigEndian.PutUint64(packed[24:32], low)
	return packed
}
