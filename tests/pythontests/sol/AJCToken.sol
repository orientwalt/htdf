// SPDX-License-Identifier: GPL-3.0
// upgrade to solidity 0.8.0
pragma solidity ^0.8.0;


contract AJCToken {

  //  no need use safemath in solidity v0.8.x
  //using SafeMath for uint256;
  
  event Transfer(address indexed from, address indexed to, uint256 value);
  event Approval(address indexed owner, address indexed spender, uint256 value);

  mapping(address => uint256) balances;
  mapping (address => mapping (address => uint256)) allowed;  

  string public constant name = "AJC chain";
  string public constant symbol = "AJC";
  uint8 public constant decimals = 18;

  constructor() {
    totalSupply = 199000000 * 10**18;
    balances[msg.sender] = totalSupply;
  }
  
  uint256 public totalSupply;

  function balanceOf(address _owner) public view returns (uint256 balance) {
    return balances[_owner];
  }

  function transfer(address _to, uint256 _value) public returns (bool) {
    require(_to != address(0));

    balances[msg.sender] = balances[msg.sender] - _value;
    balances[_to] = balances[_to] + _value;
    emit Transfer(msg.sender, _to, _value);
    return true;
  }


  function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
    uint256 _allowance = allowed[_from][msg.sender];
    require(_to != address(0));
    require (_value <= _allowance);
    balances[_from] = balances[_from] - _value;
    balances[_to] = balances[_to] + _value;
    allowed[_from][msg.sender] = _allowance - _value;
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

}
