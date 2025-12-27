// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// Interface for mintable tokens (like KawaiToken)
interface IMintableToken {
    function mint(address to, uint256 amount) external;
}

/**
 * @title MerkleDistributor
 * @dev Gas-efficient reward distribution using Merkle proofs.
 *      Supports two modes:
 *      - Mint Mode: For KAWAI rewards (mints new tokens on claim)
 *      - Transfer Mode: For USDT dividends (transfers from pre-funded balance)
 */
contract MerkleDistributor is Ownable {
    using SafeERC20 for IERC20;

    IERC20 public immutable token;
    bytes32 public merkleRoot;

    // If true, tokens are minted on claim. If false, tokens are transferred.
    bool public immutable mintOnClaim;

    // This is a packed array of booleans for gas-efficient claim tracking.
    mapping(uint256 => uint256) private claimedBitMap;

    event Claimed(uint256 index, address account, uint256 amount);
    event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot);

    /**
     * @param token_ Address of the token to distribute
     * @param mintOnClaim_ If true, mints tokens on claim (requires MINTER_ROLE).
     *                     If false, transfers tokens from contract balance.
     */
    constructor(address token_, bool mintOnClaim_) Ownable(msg.sender) {
        token = IERC20(token_);
        mintOnClaim = mintOnClaim_;
    }

    function isClaimed(uint256 index) public view returns (bool) {
        uint256 claimedWordIndex = index / 256;
        uint256 claimedBitIndex = index % 256;
        uint256 claimedWord = claimedBitMap[claimedWordIndex];
        uint256 mask = (1 << claimedBitIndex);
        return claimedWord & mask == mask;
    }

    function _setClaimed(uint256 index) private {
        uint256 claimedWordIndex = index / 256;
        uint256 claimedBitIndex = index % 256;
        claimedBitMap[claimedWordIndex] =
            claimedBitMap[claimedWordIndex] |
            (1 << claimedBitIndex);
    }

    /**
     * @notice Claim tokens using a Merkle proof.
     * @dev Caller pays gas. Tokens are either minted or transferred based on mode.
     */
    function claim(
        uint256 index,
        address account,
        uint256 amount,
        bytes32[] calldata merkleProof
    ) external {
        require(!isClaimed(index), "MerkleDistributor: Drop already claimed.");

        // Verify the merkle proof.
        bytes32 node = keccak256(abi.encodePacked(index, account, amount));
        require(
            MerkleProof.verify(merkleProof, merkleRoot, node),
            "MerkleDistributor: Invalid proof."
        );

        // Mark it as claimed
        _setClaimed(index);

        // Distribute tokens based on mode
        if (mintOnClaim) {
            // Mint new tokens directly to claimant (gas paid by claimant)
            IMintableToken(address(token)).mint(account, amount);
        } else {
            // Transfer from contract's pre-funded balance
            token.safeTransfer(account, amount);
        }

        emit Claimed(index, account, amount);
    }

    function setMerkleRoot(bytes32 _merkleRoot) external onlyOwner {
        emit MerkleRootUpdated(merkleRoot, _merkleRoot);
        merkleRoot = _merkleRoot;
    }
}
