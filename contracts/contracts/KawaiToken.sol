// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title KawaiToken
 * @notice Official KAWAI token for Kawai AI platform
 * @dev ERC20 token with minting capability and max supply cap
 * 
 * Website: https://getkawai.com
 * Docs: https://getkawai.com/docs
 * 
 * Features:
 * - Max supply: 1 billion tokens
 * - Fair launch: No initial mint
 * - Mintable by authorized distributors
 * - Burnable by token holders
 */
contract KawaiToken is ERC20, ERC20Burnable, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    // Max Supply: 1 Billion Tokens (18 decimals)
    uint256 public constant MAX_SUPPLY = 1000000000 * 10 ** 18;

    constructor(
        address defaultAdmin,
        address minter
    ) ERC20("Kawai Token", "KAWAI") {
        _grantRole(DEFAULT_ADMIN_ROLE, defaultAdmin);
        _grantRole(MINTER_ROLE, minter);
        // Fair Launch: No Initial Mint. Supply starts at 0.
    }

    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        require(totalSupply() + amount <= MAX_SUPPLY, "Max supply exceeded");
        _mint(to, amount);
    }
}
