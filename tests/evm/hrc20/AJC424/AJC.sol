pragma solidity ^0.4.24;

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

contract AJCToken {
  using SafeMath for uint256;
  
  event Transfer(address indexed from, address indexed to, uint256 value);
  event Approval(address indexed owner, address indexed spender, uint256 value);

  mapping(address => uint256) balances;
  mapping (address => mapping (address => uint256)) allowed;  

  string public constant name = "AJC chain";
  string public constant symbol = "AJC";
  uint8 public constant decimals = 18;
  address public founder = address(0);  

  constructor() public {
    totalSupply = 199000000 * 10**18;
    balances[msg.sender] = totalSupply;
    founder = msg.sender;
  }
  
  uint256 public totalSupply;

  function balanceOf(address _owner) public view returns (uint256 balance) {
    return balances[_owner];
  }

  function transfer(address _to, uint256 _value) public {
    require(_to != address(0),"Can't transfer to 0 addr");

    balances[msg.sender] = balances[msg.sender].sub(_value);
    balances[_to] = balances[_to].add(_value);
    emit Transfer(msg.sender, _to, _value);
  }


  function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
    uint256 _allowance = allowed[_from][msg.sender];
    require(_to != address(0),"Can't transfer to 0 addr");
    require (_value <= _allowance);
    balances[_from] = balances[_from].sub(_value);
    balances[_to] = balances[_to].add(_value);
    allowed[_from][msg.sender] = _allowance.sub(_value);
    emit Transfer(_from, _to, _value);
    return true;
  }


  function approve(address _spender, uint256 _value) public returns (bool) {
    require((_value == 0) || (allowed[msg.sender][_spender] == 0));
    allowed[msg.sender][_spender] = _value;
    emit Approval(msg.sender, _spender, _value);
    return true;
  }


  function allowance(address _owner, address _spender) public view returns (uint256 remaining) {
    return allowed[_owner][_spender];
  }  
  
  function changeFounder(address newFounder) public {
	  require(msg.sender == founder,"Can't transfer to 0 addr");
	  founder = newFounder;
  }  
}
