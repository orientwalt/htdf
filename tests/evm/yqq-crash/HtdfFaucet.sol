pragma solidity ^0.4.20;



contract HtdfFaucet {
    
    uint public onceAmount  = 100000000;
    address public owner ;
    
    event SendHtdf(address indexed toAddress, uint indexed amount);
    event Deposit(address indexed fromAddress, uint indexed amount);
    event SetOnceAmount(address indexed fromAddress, uint indexed amount);
    mapping (address => uint) sendRecords;
    
    function HtdfFaucet() public payable{
        owner = msg.sender;
    }
    
    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }
    
    function setOnceAmount(uint amount) public onlyOwner {
        onceAmount = amount;
        SetOnceAmount(msg.sender, amount);
    }
    
    function getOneHtdf() public {
        require( sendRecords[msg.sender] == 0 || 
            (sendRecords[msg.sender] > 0 &&  now - sendRecords[msg.sender] > 1 minutes ));
            
        require(address(this).balance >= onceAmount);
        
        msg.sender.transfer( onceAmount );
        sendRecords[msg.sender] = now;
        SendHtdf(msg.sender, onceAmount);
    }
    
    function deposit() public payable {
        Deposit(msg.sender, msg.value);
    }
    
    function() public payable{
        
    }
    
}