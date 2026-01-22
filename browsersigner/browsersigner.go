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
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/san-lab/go4337/entrypoint"
	"github.com/san-lab/go4337/signer"
	"github.com/san-lab/go4337/state"

	//. "github.com/san-lab/go4337/ui/common"
	"github.com/san-lab/go4337/userop"
)

func Init() {
	state.Register(SType, AddBrowserSigner, Unmarshal)
	aBrowserSigner := new(BrowserSigner)
	state.AddSigner(aBrowserSigner)
}

type BrowserSigner struct {
}

const SType = "Browser_Plugin"

/*
type Signer interface {
	SignMessage([]byte) ([]byte, error)
	SignEIP712([]byte) ([]byte, error)
	SignHash([]byte) ([]byte, error) //without any decorations
	Name() string
	String() string //Details
	Type() string
	Marshal() ([]byte, error)
	GetKey() any
}*/

func (bs *BrowserSigner) SignMessage(message []byte) ([]byte, error) {
	bss, err := GetBrowserSigner()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	// Set up a channel to listen for OS signals (like Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		// Wait for the signal
		<-sig
		fmt.Println("\n[Interrupt] Ctrl+C detected. Cancelling signature request...")
		// Cancel the context, which will unblock the select in SignHashWithPersonalSign
		cancel()
	}()

	return bss.SignMessage(ctx, message)
}

func (bs *BrowserSigner) SignHash([]byte) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (bs *BrowserSigner) SignEIP712(uop *userop.UserOperation, chainId *big.Int, entrypoint common.Address) ([]byte, error) {
	//return SignEIP712Way(uop, chainId, entrypoint, nil)
	bss, err := GetBrowserSigner()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	// Set up a channel to listen for OS signals (like Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		// Wait for the signal
		<-sig
		fmt.Println("\n[Interrupt] Ctrl+C detected. Cancelling signature request...")
		// Cancel the context, which will unblock the select in SignHashWithPersonalSign
		cancel()
	}()

	return bss.SignEIP712Way(ctx, uop, chainId, entrypoint, nil)

}

func (bs *BrowserSigner) SignUserop(uop *userop.UserOperation, chainId *big.Int, entryPoint common.Address) ([]byte, error) {

	switch entryPoint.Hex() {
	case entrypoint.E8Address.Hex():
		return bs.SignEIP712(uop, chainId, entryPoint)

	case entrypoint.E7Address.Hex():
		hash, err := userop.GetUserOpHashV7(uop.Pack(), entryPoint, chainId)
		if err != nil {
			return nil, err
		}
		return bs.SignMessage(hash[:])

	default:
		hash, err := userop.GetUserOpHashV6(uop, entryPoint, chainId)
		if err != nil {
			return nil, err
		}
		return bs.SignMessage(hash[:])

	}
}

func (bs *BrowserSigner) Name() string {
	return "Browser plugin"
}

func (bs *BrowserSigner) Type() string {
	return SType
}

func (bs *BrowserSigner) GetKey() any {
	return nil
}

func (bs *BrowserSigner) String() string {
	return "Browser plugin signer"
}

func (bs *BrowserSigner) Marshal() ([]byte, error) {
	return []byte{}, nil // fmt.Errorf("not implemented/needed")
}

var nameCounter = 0

func Unmarshal(bts []byte) (signer.Signer, error) {
	/*
		nameAndAddress := string(bts)
		var name, hexaddr string
		terms := strings.Split(nameAndAddress, ":")
		if len(terms) == 1 {
			name = "unnamed" + fmt.Sprint(nameCounter)
			nameCounter++
			hexaddr = terms[0]
		} else {
			name = terms[0]
			hexaddr = terms[1]
		}

		bs := new(BrowserSigner)

		bs.name = name
		bs.SignerAddress = common.HexToAddress(hexaddr)
		return bs, nil
	*/
	return nil, fmt.Errorf("Not implemented/needed")
}

// UI
func AddBrowserSigner() (err error) {
	/*
		nit := &Item{Label: "Signer name"}
		err = InputNewStringUI(nit)
		if err != nil {
			return
		}
		name, ok := nit.Value.(string)
		if !ok {
			return fmt.Errorf("invalid input for signer name")
		}

		it := &Item{Label: "Input Browser Signer Address"}
		err = InputBytes(it, -1)
		if err != nil {
			return
		}
		bt, ok := it.Value.([]byte)
		if !ok || len(bt) == 0 {
			return fmt.Errorf("invalid input")
		}

		brs := new(BrowserSigner)
		brs.name = name
		brs.SignerAddress = common.BytesToAddress(bt)
		state.AddSigner(brs)
		fmt.Println("Added browser signer", brs.Name(), brs.String())
		return
	*/
	fmt.Println("Browser Signer is a singleton")
	return
}

func (bs *BrowserSignerServer) SignEIP712Way(
	ctx context.Context,
	userOp *userop.UserOperation,
	chainId *big.Int,
	entrypoint common.Address,
	eip7702Delegate *common.Address,
) ([]byte, error) {
	// 1. Generate EIP-712 Payload using the EIP-7702 logic
	payloadJSON, err := PrepareEIP712Payload(userOp, entrypoint, chainId, eip7702Delegate)
	if err != nil {
		return nil, fmt.Errorf("failed to pack UserOp: %v", err)
	}

	// 2. Setup Request State (UUID and Channel)
	requestID := uuid.New().String() // Assuming UUID library is imported
	sigChan := make(chan []byte, 1)  // Buffered channel

	bs.mu.Lock()
	bs.requests[requestID] = &SignRequest{
		UserOp:      userOp,
		EntryPoint:  entrypoint,
		ChainID:     chainId,
		Delegate:    eip7702Delegate,
		PayloadJSON: payloadJSON,
		SigChan:     sigChan,
	}
	bs.mu.Unlock()

	// Ensure cleanup happens after signature is received or timeout
	defer func() {
		bs.mu.Lock()
		delete(bs.requests, requestID)
		bs.mu.Unlock()
	}()

	// 3. Open Browser
	url := fmt.Sprintf("%s/sign712?id=%s", bs.hostAndPort, requestID)
	//url := fmt.Sprintf("%s/sign712", bs.hostAndPort)
	fmt.Println(url)

	fmt.Printf("Listening on %s\n", bs.hostAndPort)
	if err := openBrowser(url); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please open %s manually\n", url)
	}

	fmt.Println("Waiting for signature from Browser...")

	// 4. Wait for signature
	select {
	case signatureBytes := <-sigChan:
		return signatureBytes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Minute): // Optional timeout
		return nil, fmt.Errorf("timed out waiting for signature")
	}
}

func (bs *BrowserSignerServer) SignMessage(
	ctx context.Context,
	message []byte,
) ([]byte, error) {
	// 1. Generate EIP-712 Payload using the EIP-7702 logic
	payloadJSON := "0x" + hex.EncodeToString(message)

	// 2. Setup Request State (UUID and Channel)
	requestID := uuid.New().String() // Assuming UUID library is imported
	sigChan := make(chan []byte, 1)  // Buffered channel

	bs.mu.Lock()
	bs.requests[requestID] = &SignRequest{
		PayloadJSON: payloadJSON,
		SigChan:     sigChan,
	}
	bs.mu.Unlock()

	// Ensure cleanup happens after signature is received or timeout
	defer func() {
		bs.mu.Lock()
		delete(bs.requests, requestID)
		bs.mu.Unlock()
	}()

	// 3. Open Browser
	url := fmt.Sprintf("%s/signmessage?id=%s", bs.hostAndPort, requestID)
	fmt.Println(url)

	fmt.Printf("Listening on %s\n", bs.hostAndPort)
	if err := openBrowser(url); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please open %s manually\n", url)
	}

	fmt.Println("Waiting for signature from Browser...")

	// 4. Wait for signature
	select {
	case signatureBytes := <-sigChan:
		return signatureBytes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Minute): // Optional timeout
		return nil, fmt.Errorf("timed out waiting for signature")
	}
}

// ---------------BrowserSignerServer------------------
var (
	// Private instance variable to hold the singleton
	browserSignerInstance *BrowserSignerServer
	// Primitives to ensure thread-safe, one-time initialization
	once sync.Once
)

// GetBrowserSigner safely initializes and returns the single instance.
func GetBrowserSigner() (*BrowserSignerServer, error) {
	fmt.Println("Getting browser signer")
	var initErr error

	// This code block runs only once, across all goroutines
	once.Do(func() {
		// Call the constructor that starts the server
		instance, err := NewBrowserSignerServer()
		if err != nil {
			initErr = fmt.Errorf("failed to start BrowserSigner: %w", err)
			return
		}
		browserSignerInstance = instance
	})

	// Return the result of the initialization
	if initErr != nil {
		return nil, initErr
	}
	return browserSignerInstance, nil
}

type SignRequest struct {
	UserOp      *userop.UserOperation
	EntryPoint  common.Address
	ChainID     *big.Int
	Delegate    *common.Address // EIP-7702 Delegate Address
	PayloadJSON string
	SigChan     chan []byte // Channel to send the signature back
}

// BrowserSigner manages the long-lived HTTP server and active requests.
type BrowserSignerServer struct {
	server      *http.Server
	listener    net.Listener
	requests    map[string]*SignRequest // Key: Request ID (UUID)
	mu          sync.Mutex              // Protects access to the requests map
	hostAndPort string
}

var signerport = 8761

// NewBrowserSigner creates and starts the long-running HTTP server.
func NewBrowserSignerServer() (*BrowserSignerServer, error) {
	mux := http.NewServeMux()
	bs := &BrowserSignerServer{
		requests: make(map[string]*SignRequest),
	}

	// Bind to port 0 to get a random free port
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", signerport))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on local port: %v", err)
	}
	bs.listener = listener
	bs.hostAndPort = fmt.Sprintf("http://localhost:%v", signerport)

	bs.server = &http.Server{
		Handler: mux,
	}

	// Register handlers using methods on the struct
	mux.HandleFunc("/sign712", bs.sign712Handler)
	mux.HandleFunc("/signmessage", bs.signMessageHandler)
	mux.HandleFunc("/submit", bs.submitHandler)

	// Start Server in a persistent Goroutine
	go func() {
		if err := bs.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("BrowserSigner server error: %v", err)
		}
	}()

	return bs, nil
}

// Shutdown gracefully shuts down the HTTP server.
func (bs *BrowserSignerServer) Shutdown(ctx context.Context) error {
	return bs.server.Shutdown(ctx)
}

func (bs *BrowserSignerServer) signMessageHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get Request ID and Payload
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	bs.mu.Lock()
	req, ok := bs.requests[requestID]
	bs.mu.Unlock()

	if !ok {
		http.Error(w, "Invalid or expired request ID", http.StatusNotFound)
		return
	}

	// 2. Serve Template with Payload and Request ID
	tmpl, err := template.New("signMessage").Parse(signMessagePageTemplate)
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	// The template must be updated to use the request ID
	tmpl.Execute(w, struct{ Payload, RequestID string }{
		Payload:   req.PayloadJSON,
		RequestID: requestID,
	})
}

func (bs *BrowserSignerServer) sign712Handler(w http.ResponseWriter, r *http.Request) {
	// 1. Get Request ID and Payload
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	bs.mu.Lock()
	req, ok := bs.requests[requestID]
	bs.mu.Unlock()

	if !ok {
		http.Error(w, "Invalid or expired request ID", http.StatusNotFound)
		return
	}

	// 2. Serve Template with Payload and Request ID
	tmpl, err := template.New("sign712").Parse(signPageTemplate)
	if err != nil {
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	// The template must be updated to use the request ID
	tmpl.Execute(w, struct{ Payload, RequestID string }{
		Payload:   req.PayloadJSON,
		RequestID: requestID,
	})
}

func (bs *BrowserSignerServer) submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var result struct {
		Signature string `json:"signature"`
		RequestID string `json:"requestID"` // New field
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	bs.mu.Lock()
	req, ok := bs.requests[result.RequestID]
	bs.mu.Unlock()

	if !ok {
		http.Error(w, "Invalid or expired request ID", http.StatusNotFound)
		return
	}

	// Process signature and send back
	signatureHex := strings.TrimPrefix(result.Signature, "0x")
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		http.Error(w, "Invalid signature format", http.StatusBadRequest)
		return
	}

	// Send signature to the waiting goroutine
	req.SigChan <- signatureBytes

	w.WriteHeader(http.StatusOK)
}
