// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

contract Son {
    uint public x;
    constructor(uint a) payable {
        x = a;
    }
}

contract Create {
    Son son1 = new Son(4); // will be executed as part of Create's constructor
    Son son2 = new Son(4); // will be executed as part of Create's constructor
    Son son3 = new Son(4); // will be executed as part of Create's constructor
    Son son4 = new Son(4); // will be executed as part of Create's constructor

    function createSon(uint arg) public payable {
        Son newSon = new Son(arg);
        require( newSon.x() == arg, "===NOT EQUAL===");
    }


    function createSonEx(uint arg) public payable {
        for(uint i = 0; i < arg; i++) {
            Son newSon = new Son(arg);
            require( newSon.x() == arg, "===NOT EQUAL===");
        }
    }


}