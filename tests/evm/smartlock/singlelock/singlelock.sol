pragma solidity ^0.4.20;

// inferring block is more secure than inferring contract
// to hijack the rpc traffic line is much easier.
// |-block
// |---txs
//   |---From, To(contract Addr), Data, Log(True) 
// unlocking signal should be dynamic to prevent from brute-forcing.
contract SingleLock {
    address private key = address(0);
    bool public unlocked = false;
    uint constant timeout = 60 * 10;
    //
    function SingleLock() public {
        key = msg.sender;
    }
    function Lock() public {
        unlocked = false;
    }
    // slock unlocks itself.
    bytes32 public passcode;
    uint public valid = 0;
    function Remaining() public {
        return valid - now;
    }
    function UpdatePass(uint pass) public {
        require(users[msg.sender]);
        passcode = sha256(pass);//keccak256,ripemd160
        valid = now + timeout; 
    }
    function Unlock(uint pass) public {
        require(users[msg.sender]);
        require(Remaining()>0);
        require(sha256(pass)==passcode);
        unlocked = true;
    }

    // app unlocks the slock.
    function AppUnlock() public returns (bool) {
        require(users[msg.sender]);
        unlocked = true;
        return true;
    }
    // change admin
    function changeKey(address addr) public {
	  require(msg.sender == key);
	  key = addr;
    }  
}