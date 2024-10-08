A Golang tool to fabricate and manipulate ERC 4337 UserOperations for Ethereum blockchains

v0.1.0:  The first version that actually lets you import your ABI's, define your UserOperation, sign it (ECDSA, imported private key only atm), and export the UserOperation as a Remix-compatible tuple / Bundler standard JSON.
Still a lot to do...

FastForward to
v0.2.3: ·rd Party APIs: Paymasters and Bundlers. Alechemy API reasonably well integrated as well.
It helps if you have and API_KEY and "paymaster policy" from Alchemy/Stackup/Biconomy/Pimlico/... (tbc)


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
WhitelistingPaymaster (V6) | 0x6f26BCA296546533f3653360f8C501FDAf2cCA72
SimpleAccountFactory (V7) | [0x6Ac14e603A2742fB919248D66c8ecB05D8Aec1e9](https://sepolia.etherscan.io/address/0x6Ac14e603A2742fB919248D66c8ecB05D8Aec1e9)
SimpleAccountFactory (V6) | 0x9406Cc6185a346906296840746125a0E44976454
Greeter (sample contract) | [0x5Fc1328D1bA215852Ad445510Dc72fFB29718C33](https://sepolia.etherscan.io/address/0x5Fc1328D1bA215852Ad445510Dc72fFB29718C33)
Universal Factory | [0xc0721E8bcF8Ede76e5b046b35B435B3C2B3303b9](https://sepolia.etherscan.io/address/0xc0721E8bcF8Ede76e5b046b35B435B3C2B3303b9) 

## Holesky deployment

Contract | Address
------------- | -------------
Entrypoint (V7) | 0x0000000071727De22E5E9d8BAf0edAc6f37da032
SimpleAccountFactory (V7) | 0xd8c9e8F357810A5E702764F626d4A380Ca462FB1
WhitelistingPaymaster (V7) | 0xAa3f81F9b2F12a36037E06c7e009d282546B6249
Deployer (aka Univ Factory) | 0x2b1048241475CD7033515f28653362202467953a
Entrypoint (V6) | 0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789
SimpleAccountFactory (V6) | 0x79Ef5e7e18A5719D0b8aF4e7EA0A54818441aA06
WhitelistingPaymasterV6 | 0xc0CB75eD4Ce275f4A3dE293aa5f98b15FEb5907c

