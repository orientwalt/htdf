// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;


contract Constructor {

    uint256 public onceAmount = 100000;
    constructor(uint256 amount) {
        onceAmount = amount;
    }

    function setAmount(uint256 amount) public {
        onceAmount = amount;
    }
}