// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title MockStablecoin
 * @dev Simple ERC20 stablecoin for testing on testnet
 *      Mimics USDC/USDT behavior with 6 decimals
 */
contract MockStablecoin is ERC20 {
    constructor() ERC20("Mock Stablecoin", "MOCK") {
        // Mint 1 million tokens to deployer for testing
        _mint(msg.sender, 1000000 * 10 ** decimals());
    }

    // Allow anyone to mint for testing purposes
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    // Override decimals to match real stablecoins (6 decimals)
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
