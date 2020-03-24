pragma solidity ^0.6.0;

interface IERC20 {
	// 查询发行量
	function totalSupply() external view returns (uint256);
	// 查询余额
	function balanceOf(address who) external view returns (uint256);
	// 授权余额查询
	function allowance(address owner, address spender) external view returns (uint256);
	// 转账
	function transfer(address to, uint256 value) external returns (bool);
	// 授权
	function approve(address spender, uint256 value) external returns (bool);
	// 利用授权转账
	function transferFrom(address from, address to, uint256 value) external returns (bool);
 	// Transfer事件
	event Transfer(address indexed from, address indexed to,uint256 value);
	// approve事件
	event Approval(address indexed owner, address indexed spender, uint256 value);
}