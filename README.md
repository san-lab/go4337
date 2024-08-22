A Golang tool to fabricate and manipulate ERC 4337 UserOperations for Ethereum blockchains

v0.1.0:  The first version that actually lets you import your ABI's, define your UserOperation, sign it (ECDSA, imported private key only atm), and export the UserOperation as a Remix-compatible tuple / Bundler standard JSON.
Still a lot to do...

1. Create a UserOperation
2. Set all the bytecode/ABI fields by selecting smart contract methods form ABI and setting their arguments values
3. Sign it (ECDSA, as implemented in EthInfinitism’s SimpleWallet)
4. Prepare PaymasterData complete with the signature (again, ECDSA, as implemented in EthInfinitism’s VerifyingPaymaster)
5. Export the UserOperation as 
    - V6 Bundler json
    - V6 Remix tuple
    - V7 Bundler json
    - V7 Remix tuple
    - cURL call to Alchemy Bundler

---

## Sepolia deployment

Contract | Address
------------- | -------------
v7 Entrypoint | [0x0000000071727De22E5E9d8BAf0edAc6f37da032](https://sepolia.etherscan.io/address/0x0000000071727De22E5E9d8BAf0edAc6f37da032)
WhitelistingPaymaster (V7) | [0xd8c9e8F357810A5E702764F626d4A380Ca462FB1](https://sepolia.etherscan.io/address/0xd8c9e8F357810A5E702764F626d4A380Ca462FB1)
SimpleAccountFactory (V7) | [0x6Ac14e603A2742fB919248D66c8ecB05D8Aec1e9](https://sepolia.etherscan.io/address/0x6Ac14e603A2742fB919248D66c8ecB05D8Aec1e9)
Greeter (sample contract) | [0x5Fc1328D1bA215852Ad445510Dc72fFB29718C33](https://sepolia.etherscan.io/address/0x5Fc1328D1bA215852Ad445510Dc72fFB29718C33)
Universal Factory | [0xc0721E8bcF8Ede76e5b046b35B435B3C2B3303b9](https://sepolia.etherscan.io/address/0xc0721E8bcF8Ede76e5b046b35B435B3C2B3303b9) 

## Holesky deployment

Contract | Address
------------- | -------------
SimpleAccountFactory (V7) | 0xd8c9e8F357810A5E702764F626d4A380Ca462FB1
WhitelistingPaymaster (V7) | 0xAa3f81F9b2F12a36037E06c7e009d282546B6249
Deployer (aka Univ Factory) | 0x2b1048241475CD7033515f28653362202467953a

