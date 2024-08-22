// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Deployer {
    event ContractDeployed(address indexed contractAddress);

    // Function to deploy a contract with arbitrary bytecode
    function deploy(bytes memory bytecode) public returns (address) {
        address addr;
        assembly {
            // Deploy using CREATE opcode with `create` function
            // `0` specifies the amount of Ether to send (here, none)
            // `add(bytecode, 0x20)` is the offset where the bytecode starts in memory
            // `mload(bytecode)` loads the size of the bytecode
            addr := create(0, add(bytecode, 0x20), mload(bytecode))
        }

        require(addr != address(0), "Deployment failed");
        emit ContractDeployed(addr);
        return addr;
    }
}
