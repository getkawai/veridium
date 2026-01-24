// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Capped.sol";
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
 * - Max supply: 1 billion tokens (enforced by ERC20Capped)
 * - Fair launch: No initial mint
 * - Mintable by authorized distributors
 * - Burnable by token holders
 */
contract KawaiToken is ERC20, ERC20Burnable, ERC20Capped, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    constructor(
        address defaultAdmin,
        address minter
    ) ERC20("Kawai Token", "KAWAI") ERC20Capped(1000000000 * 10 ** 18) {
        _grantRole(DEFAULT_ADMIN_ROLE, defaultAdmin);
        _grantRole(MINTER_ROLE, minter);
        // Fair Launch: No Initial Mint. Supply starts at 0.
    }

    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        _mint(to, amount); // ERC20Capped automatically checks cap
    }

    /**
     * @dev Override required by Solidity for multiple inheritance
     */
    function _update(address from, address to, uint256 value) 
        internal 
        override(ERC20, ERC20Capped) 
    {
        super._update(from, to, value);
    }
}
