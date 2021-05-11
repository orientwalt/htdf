// SPDX-License-Identifier: GPL-3.0
// yqq 2020-12-11
// test contract 

// upgrade to v0.8.0 , yqq,2021-05-11
pragma solidity ^0.8.0;


contract HtdfFaucet {
    
    uint256 public onceAmount;
    address public owner ;
    
    event SendHtdf(address indexed toAddress, uint256 indexed amount);
    event Deposit(address indexed fromAddress, uint256 indexed amount);
    event SetOnceAmount(address indexed fromAddress, uint256 indexed amount);
    mapping (address => uint256) sendRecords;
    
    constructor() payable {
        onceAmount = 100000000; // once 1 HTDF
        owner = msg.sender;
    }
    
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
    
    function setOnceAmount(uint256 amount) public onlyOwner {
        onceAmount = amount;
        emit SetOnceAmount(msg.sender, amount);
    }
    
    function getOneHtdf() public {
        require( sendRecords[msg.sender] == 0 || 
            (sendRecords[msg.sender] > 0 &&  block.timestamp - sendRecords[msg.sender] > 1 minutes ));
            
        require(address(this).balance >= onceAmount);

        //  update time before transfer to against re-entrancy attack
        sendRecords[msg.sender] = block.timestamp;

        // transfer only use 2300 gas, safe against re-entrancy attack
        payable(msg.sender).transfer( onceAmount );

        emit SendHtdf(msg.sender, onceAmount);
    }
    
    function deposit() public payable {
        emit Deposit(msg.sender, msg.value);
    }


    // fallback fucnction
    fallback() external payable {
        if(msg.value > 0) {
            emit Deposit(msg.sender, msg.value);
        }
    }

    receive() external payable{
        if(msg.value > 0) {
            emit Deposit(msg.sender, msg.value);
        }
    }

    // desctruct the contract, send all left coins to owner
    function kill() public onlyOwner {
         selfdestruct( payable(owner));
    }
    
}