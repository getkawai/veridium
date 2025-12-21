// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title PaymentVault
 * @dev Handles user deposits in USDT for AI service credits.
 */
contract PaymentVault is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable usdt;

    event Deposited(address indexed user, uint256 amount);
    event Withdrawn(address indexed to, uint256 amount);

    constructor(address _usdt, address initialOwner) Ownable(initialOwner) {
        require(_usdt != address(0), "Invalid USDT address");
        usdt = IERC20(_usdt);
    }

    /**
     * @notice Deposit USDT to get service credits.
     * @param _amount Amount of USDT to deposit.
     */
    function deposit(uint256 _amount) external nonReentrant {
        require(_amount > 0, "Amount must be > 0");
        usdt.safeTransferFrom(msg.sender, address(this), _amount);
        emit Deposited(msg.sender, _amount);
    }

    /**
     * @notice Withdraw USDT from the vault (Owner only - for revenue distribution).
     * @param _to Recipient address.
     * @param _amount Amount to withdraw.
     */
    function withdraw(
        address _to,
        uint256 _amount
    ) external onlyOwner nonReentrant {
        require(_to != address(0), "Invalid recipient");
        require(
            _amount <= usdt.balanceOf(address(this)),
            "Insufficient balance"
        );
        usdt.safeTransfer(_to, _amount);
        emit Withdrawn(_to, _amount);
    }
}
