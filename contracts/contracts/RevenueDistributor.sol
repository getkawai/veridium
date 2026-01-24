// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title RevenueDistributor
 * @dev Gas-efficient stablecoin revenue sharing distribution using Merkle proofs.
 *      Transfers tokens from pre-funded contract balance.
 *      Used for distributing revenue share to KAWAI token holders.
 *      Supports any ERC20 stablecoin (USDC, USDT, DAI, etc).
 */
contract RevenueDistributor is Ownable {
    using SafeERC20 for IERC20;

    IERC20 public immutable token;
    bytes32 public merkleRoot;

    // This is a packed array of booleans for gas-efficient claim tracking.
    mapping(uint256 => uint256) private claimedBitMap;

    event Claimed(uint256 index, address account, uint256 amount);
    event MerkleRootUpdated(bytes32 oldRoot, bytes32 newRoot);

    /**
     * @param token_ Address of the stablecoin token to distribute (USDC/USDT/etc)
     */
    constructor(address token_) Ownable(msg.sender) {
        token = IERC20(token_);
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
     * @notice Claim stablecoin revenue using a Merkle proof.
     * @dev Caller pays gas. Tokens are transferred from contract's pre-funded balance.
     */
    function claim(
        uint256 index,
        address account,
        uint256 amount,
        bytes32[] calldata merkleProof
    ) external {
        require(!isClaimed(index), "RevenueDistributor: Drop already claimed.");

        // Verify the merkle proof.
        bytes32 node = keccak256(abi.encodePacked(index, account, amount));
        require(
            MerkleProof.verify(merkleProof, merkleRoot, node),
            "RevenueDistributor: Invalid proof."
        );

        // Mark it as claimed
        _setClaimed(index);

        // Transfer stablecoin from contract's pre-funded balance
        token.safeTransfer(account, amount);

        emit Claimed(index, account, amount);
    }

    function setMerkleRoot(bytes32 _merkleRoot) external onlyOwner {
        emit MerkleRootUpdated(merkleRoot, _merkleRoot);
        merkleRoot = _merkleRoot;
    }
}
