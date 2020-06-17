pragma solidity ^0.4.20;


contract ERC20Token {
    uint256 public totalSupply;
    function balanceOf(address _owner) public constant returns (uint256 balance);
    function transfer(address _to, uint256 _value) public returns (bool success);
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success);
    function approve(address _spender, uint256 _value) public returns (bool success);
    function allowance(address _owner, address _spender) public constant returns (uint256 remaining);
    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}
// premise:
// all exchange users including the exchange founder need to approve their accounts for this exchange contract in the corresponding hrc20 token contract.
// checkpoints:
// like in function htdf2token, is it possible to call transferfrom without signature?
contract Exchange {
    address public founder = address(0);
    mapping(address => uint) public htdf;
    function selfincrease(uint256 amount) public returns (bool) {
        if (!msg.sender.send(amount)) {
            return false;
        }
        return true;
    }

    function htdf2token(address tokenAddr) public payable {
        // htdf: sender to contract
        htdf[msg.sender] += msg.value;
        // token: founder to sender
        ERC20Token(tokenAddr).transferFrom(founder, msg.sender, msg.value);
    }

    function token2htdf(address tokenAddr, uint256 value) public payable returns (bool) {
        // token: sender to founder
        ERC20Token(tokenAddr).transferFrom(msg.sender, founder, value);
        // htdf:  founder to sender
        if (!msg.sender.send(value)) {
            return false;
        }
        return true;
    }
}