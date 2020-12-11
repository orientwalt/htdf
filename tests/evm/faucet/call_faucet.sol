/*
yqq  2020-12-11
test contract call contract
*/

pragma solidity ^0.4.20;

import "./htdf_faucet.sol";

contract CallFaucet {

    HtdfFaucet public faucet;
    uint256 public stackDepth = 0;
    address public addr;
    address public owner;

    function CallFaucet() public payable {
        // hard coding the faucet contract address
        addr = address(0x18EDA861679664967c067bA3068414339E5B49e9);
        faucet = HtdfFaucet(addr);
        owner = msg.sender;
    }

    function  getOneHtdf() public {
        stackDepth = 0;
        faucet.getOneHtdf(); // it's ok
    }

    // fallback function
    function () external payable {
       
    }

}