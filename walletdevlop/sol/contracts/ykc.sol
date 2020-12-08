pragma solidity ^0.6.0;

import "./IERC20.sol";
import "./SafeMath.sol";

contract ERC20 is IERC20 {
  // 使用安全函数
  using SafeMath for uint256;
	// 记录用户余额的结构 
  mapping (address => uint256) private _balances;
	// 记录用户授权的结构 from => to => value
  mapping (address => mapping (address => uint256)) private _allowed;
	// 总发行量
  uint256 private _totalSupply;
  // token名
  string public symbol;
  //管理员
  address private owner;
  
  constructor(string memory _sym) public {
      symbol = _sym;
      owner = msg.sender;
  }
	// 挖矿
  function mint(address to, uint256 value) public {
      require(msg.sender == owner);
      _totalSupply = _totalSupply.add(value);
      _balances[to] = _balances[to].add(value);
      emit Transfer(address(0), to, value);
  }
	// 总发行量
  function totalSupply() override public view returns (uint256) {
    return _totalSupply;
  }
  function balanceOf(address owner) override public view returns (uint256) {
    return _balances[owner];
  }
  function approve(address spender, uint256 value) override public returns (bool) {
    require(spender != address(0));

    _allowed[msg.sender][spender] = value;
    emit Approval(msg.sender, spender, value);
    return true;
  }
  function allowance(address owner, address spender) override public view returns (uint256) {
    return _allowed[owner][spender];
  }
   function transferFrom(
    address from,
    address to,
    uint256 value
  )
  override
    public
    returns (bool)
  {
    // 用户余额充足
    require(value <= _balances[from]);
     // 用户授权余额充足
    require(value <= _allowed[from][msg.sender]);
    require(to != address(0));
		// 划账 A- B+
    _balances[from] = _balances[from].sub(value);
    _balances[to] = _balances[to].add(value);
    _allowed[from][msg.sender] = _allowed[from][msg.sender].sub(value);
    emit Transfer(from, to, value);
    return true;
  }
  function transfer(address to, uint256 value) override public returns (bool) {
    // 转出用户余额充足
  	require(value <= _balances[msg.sender]);
    require(to != address(0));
		// 调整账本
    _balances[msg.sender] = _balances[msg.sender].sub(value);
    _balances[to] = _balances[to].add(value);
    emit Transfer(msg.sender, to, value);
    return true;
  }
}