pragma solidity ^0.4.20;

library SafeMath {
  function mul(uint256 a, uint256 b) internal pure returns (uint256) {
	if (a == 0) {
		return 0;
	}
	uint256 c = a * b;
	assert(c / a == b);
	return c;
  }

  function div(uint256 a, uint256 b) internal pure returns (uint256) {
	uint256 c = a / b;
	return c;
  }

  function sub(uint256 a, uint256 b) internal pure returns (uint256) {
	assert(b <= a);
	return a - b;
  }

  function add(uint256 a, uint256 b) internal pure returns (uint256) {
	uint256 c = a + b;
	assert(c >= a);
	return c;
  }
}

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
// like in function htdf2token, is it possible to call transferfrom without signature? absolutely no.
contract Exchange {
    using SafeMath for uint256;
    address public founder = address(0);
    uint accumulated = 0;
    mapping(address => uint) public htdf;
    // constructor
    function Exchange() public {
        founder = msg.sender;
    }
    // increase exchange ether amount
    function donate() public payable{
        accumulated = accumulated.add(msg.value);
    }
    //
    function balanceOf(address addr) public view returns (uint256) {
        return address(addr).balance;
    }
    function tokenbalanceOf(address tokenAddr,address addr) public view returns (uint256) {
        return ERC20Token(tokenAddr).balanceOf(addr);
    }
    // test passed
    function htdf2token(address tokenAddr) public payable {
        // htdf: sender to contract
        accumulated = accumulated.add(msg.value);
        // token: founder to sender
        ERC20Token(tokenAddr).transferFrom(founder, msg.sender, msg.value);
    }
    // test passed
    function token2htdf(address tokenAddr, uint256 amount) public payable returns (bool) {
        require(accumulated > amount);
        // token: sender to founder
        if (!ERC20Token(tokenAddr).transferFrom(msg.sender, founder, amount)) {
            return false;
        }
        // htdf:  founder to sender
        if (!msg.sender.call.value(amount)(true, 3)) {
            return false;
        }
        accumulated = accumulated.sub(amount);
        return true;
    }
    //
    function changeFounder(address newFounder) public {
	  require(msg.sender == founder);
	  founder = newFounder;
    }  
}