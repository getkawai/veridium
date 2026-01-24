// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title PaymentVault
 * @notice Handles user deposits in stablecoin (USDC/USDT) for AI service credits
 * @dev Secure vault with ReentrancyGuard and access control
 * 
 * Website: https://getkawai.com
 * Docs: https://getkawai.com/docs
 * 
 * Features:
 * - Accepts USDC (mainnet) or MockStablecoin (testnet)
 * - Converts deposits to AI service credits
 * - Owner-only withdrawal for revenue distribution
 */
contract PaymentVault is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable stablecoin;

    event Deposited(address indexed user, uint256 amount);
    event Withdrawn(address indexed to, uint256 amount);

    constructor(address _stablecoin, address initialOwner) Ownable(initialOwner) {
        require(_stablecoin != address(0), "Invalid stablecoin address");
        stablecoin = IERC20(_stablecoin);
    }

    /**
     * @notice Deposit stablecoin to get service credits.
     * @param _amount Amount of stablecoin to deposit.
     */
    function deposit(uint256 _amount) external nonReentrant {
        require(_amount > 0, "Amount must be > 0");
        stablecoin.safeTransferFrom(msg.sender, address(this), _amount);
        emit Deposited(msg.sender, _amount);
    }

    /**
     * @notice Withdraw stablecoin from the vault (Owner only - for revenue distribution).
     * @param _to Recipient address.
     * @param _amount Amount to withdraw.
     */
    function withdraw(
        address _to,
        uint256 _amount
    ) external onlyOwner nonReentrant {
        require(_to != address(0), "Invalid recipient");
        require(
            _amount <= stablecoin.balanceOf(address(this)),
            "Insufficient balance"
        );
        stablecoin.safeTransfer(_to, _amount);
        emit Withdrawn(_to, _amount);
    }
}
