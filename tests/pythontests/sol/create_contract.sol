pragma solidity >0.4.99 <0.9.0;

contract Son {
    uint public x;
    constructor(uint a) public payable {
        x = a;
    }

    function calc(uint a, uint b) public returns(uint) {
        return a + b;
    }
}

contract Parent {
    Son son = new Son(4); // will be executed as part of C's constructor

    function createSon(uint arg) public payable {
        Son newSon = new Son(arg);
        require( newSon.x() == arg, "===NOT EQUAL===");
    }

}