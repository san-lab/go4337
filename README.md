A Golang tool to fabricate and manipulate ERC 4337 UserOperations for Ethereum blockchains

v0.1.0:  The first version that actually lets you import your ABI's, define your UserOperation, sign it (ECDSA, imported private key only atm), and export the UserOperation as a Remix-compatible tuple / Bundler standard JSON.
Still a lot to do...

With the tool, one can atm:
1)	Create a UserOperation
2)	Set all the bytecode/ABI fields by selecting smart contract methods form ABI and setting their arguments values
3)	Sign it (ECDSA, as implemented in EthInfinitism’s SimpleWallet)
4)	Prepare PaymasterData complete with the signature (again, ECDSA, as implemented in EthInfinitism’s VerifyingPaymaster)
5)	Export the UserOperation as 
a.	V6 Bundler json
b.	V6 Remix tuple
c.	V7 Bundler json
d.	V7 Remix tuple

