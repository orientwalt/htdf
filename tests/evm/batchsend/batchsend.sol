pragma solidity ^0.4.20;
//https://solidity.readthedocs.io/en/develop/units-and-global-variables.html?highlight=transfer#address-related
//https://ethereum.stackexchange.com/questions/19341/address-send-vs-address-transfer-best-practice-usage
//https://ethereum.stackexchange.com/questions/78124/is-transfer-still-safe-after-the-istanbul-update
// transfer vs. send
// tranfer returns nothing on failure while send returns false on failure

// https://medium.com/@rsripathi781/6-payable-functions-in-solidity-smartcontract-ethereum-d2535e346dc1
// Payable functions provide a mechanism to collect / receive funds in ethers to your contract . Payable functions are annotated with payable keyword.
// In the above example payme function is annotated with payable keyword, which translates to that you can send ethers to payme function.
// Where are the ethers stored in payable functions?
// All the ethers sent to payable functions are owned by contract. In the above example Sample contract owns all of the ethers.
contract  PayableTest {
    uint accumulated =0;
    address public founder = address(0);
    function receive() public payable{
        accumulated += msg.value;
    }

    // function getbalance() public view returns (uint256 balance) {
    //     return msg.sender.balance();
    // }

    // function setbalance(bool tofounder,uint amount) public {
    //     address receiver = msg.sender;
    //     if (tofounder){
    //         receiver = founder;
    //     }
    //     receiver.balance(amount);
    // }

    function transfer(bool tofounder,uint amount) public {
        address receiver = msg.sender;
        if (tofounder){
            receiver = founder;
        }
        founder.transfer(amount);
    }

    function send(bool tofounder,uint amount) public returns(bool){
        address receiver = msg.sender;
        if (tofounder){
            receiver = founder;
        }
        if (receiver.send(amount)) {
            // No need to call throw here, just reset the amount owing
            return false;
        }
        return true;
    }

    function callex(bool tofounder,uint amount) public returns(bool) {
        address receiver = msg.sender;
        if (tofounder){
            receiver = founder;
        }
        if (!receiver.call.gas(30000).value(amount)()){
            return false;
        }
        return true;
    }   
    // junying-todo, 2020-06-17, test passed
    function call(bool tofounder,uint amount) public returns(bool)  {
        address receiver = msg.sender;
        if (tofounder){
            receiver = founder;
        }
        if (!receiver.call.value(amount)(true, 3)) {
            // No need to call throw here, just reset the amount owing
            return false;
        }
        return true;
    }
}

contract BatchSend {
    address public founder = address(0);
    // only founder has privilege to call this function
    function distribute(address[] dests, uint256[] values) public payable {
        require(msg.sender == founder);
        for (uint i = 0; i < dests.length; i++) {
            dests[i].transfer(values[i]);
        }
    }
}

// this contract is designed to distribute qutoa to the parties.
// every receiver will call withdraw to receive the corresponding quota.
contract QuotaReceive {
    // quota mapp, which is initialized when contract creating.
    mapping(address => uint) public qutoas;
    // constructor
    function QuotaReceive(
        address[] dests,
        uint256[] values
    ) public {
        for (uint i = 0; i < dests.length; i++) {
            qutoas[dests[i]]=values[i];
        }
    }
    // qutoa receiving
    function withdraw() public returns (bool) {
        uint amount = qutoas[msg.sender];
        if (amount > 0) {
            // It is important to set this to zero because the recipient
            // can call this function again as part of the receiving call
            // before `send` returns.
            qutoas[msg.sender] = 0;

            if (!msg.sender.send(amount)) {
                // No need to call throw here, just reset the amount owing
                qutoas[msg.sender] = amount;
                return false;
            }
        }
        return true;
    }
    // getbalance
    function balanceOf(address _owner) public view returns (uint256 balance) {
        return qutoas[_owner];
    }
}