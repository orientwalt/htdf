/*
	The following is an extremely basic example of a solidity contract.
	It takes a string upon creation and then repeats it when greet() is called.
*/

pragma solidity ^0.8.0;

/// @title Greeter
/// @author Cyrus Adkisson
// The contract definition. A constructor of the same name will be automatically called on contract creation.
contract Suicide {

    // At first, an empty "address"-type variable of the name "creator". Will be set in the constructor.
    address creator;

    // The constructor. It accepts a string input and saves it to the contract's "greeting" variable.
    constructor() payable {
        creator = msg.sender;
    }

     /**********
     Standard kill() function to recover funds
     **********/
    function kill() public payable {
        if (msg.sender == creator){
            selfdestruct( payable(msg.sender));
        }
    }


    receive() external payable{

    }

    fallback() external payable{

    }

}