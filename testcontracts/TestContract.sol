pragma solidity ^0.8.6;

contract TestContract {
    uint64 private sequence;
    string private message;
    mapping(uint64 => address) public Origin;
    mapping(uint64 => address) public Sender;
    mapping(uint64 => string) public Message;

    function Sequence() public view returns (uint64) {
        return sequence;
    }

    function setMessage(string calldata _msg) public {
        Origin[sequence] = tx.origin;
        Sender[sequence] = msg.sender;
        Message[sequence] = _msg;
        message = _msg;
        sequence++;
        return;
    }

    function CurrentMessage() public view returns (string memory) {
        return message;
    }
}
