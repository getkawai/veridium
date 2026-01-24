// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

// Interface for mintable tokens (like KawaiToken)
interface IMintableToken {
    function mint(address to, uint256 amount) external;
}

/**
 * @title ReferralRewardDistributor
 * @notice Specialized contract for distributing referral rewards (KAWAI tokens only)
 * @dev Uses Merkle proofs for gas-efficient batch distribution
 * 
 * Website: https://getkawai.com
 * Docs: https://getkawai.com/docs
 * 
 * Note: USDT rewards are off-chain credits only (stored in KV, not on-chain)
 * 
 * Features:
 * - Period-based Merkle settlements
 * - Batch claiming for multiple periods
 * - Tracks unique referrers
 * - Emergency pause mechanism
 */
contract ReferralRewardDistributor is Ownable, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    IERC20 public immutable kawaiToken;
    
    bytes32 public merkleRoot;
    uint256 public currentPeriod;
    
    // Track claimed rewards per period per user
    // period => user => claimed
    mapping(uint256 => mapping(address => bool)) public hasClaimed;
    
    // Track unique referrers (for totalReferrers count)
    mapping(address => bool) public hasClaimedAnyPeriod;
    
    // Store merkle roots per period for historical claims
    mapping(uint256 => bytes32) public periodMerkleRoots;
    
    // Track total rewards distributed
    uint256 public totalKawaiDistributed;
    uint256 public totalReferrers;

    event RewardsClaimed(
        uint256 indexed period,
        address indexed user,
        uint256 kawaiAmount
    );
    event MerkleRootUpdated(
        uint256 indexed period,
        bytes32 oldRoot,
        bytes32 newRoot
    );
    event PeriodAdvanced(uint256 oldPeriod, uint256 newPeriod);

    /**
     * @param kawaiToken_ Address of KAWAI token (mintable)
     */
    constructor(
        address kawaiToken_
    ) Ownable(msg.sender) {
        require(kawaiToken_ != address(0), "Invalid KAWAI address");
        kawaiToken = IERC20(kawaiToken_);
        currentPeriod = 1;
    }

    /**
     * @notice Claim referral rewards using Merkle proof
     * @param period Period number
     * @param kawaiAmount Amount of KAWAI tokens to claim
     * @param merkleProof Merkle proof for verification
     */
    function claimRewards(
        uint256 period,
        uint256 kawaiAmount,
        bytes32[] calldata merkleProof
    ) external nonReentrant whenNotPaused {
        require(period <= currentPeriod, "Invalid period");
        require(!hasClaimed[period][msg.sender], "Already claimed for this period");
        require(kawaiAmount > 0, "No rewards to claim");

        // Verify Merkle proof using period-specific root
        bytes32 leaf = keccak256(
            abi.encodePacked(period, msg.sender, kawaiAmount)
        );
        bytes32 periodRoot = periodMerkleRoots[period];
        require(periodRoot != bytes32(0), "Period not settled");
        require(
            MerkleProof.verify(merkleProof, periodRoot, leaf),
            "Invalid proof"
        );

        // Mark as claimed
        hasClaimed[period][msg.sender] = true;

        // Mint KAWAI tokens
        IMintableToken(address(kawaiToken)).mint(msg.sender, kawaiAmount);
        totalKawaiDistributed += kawaiAmount;

        // Only increment if first claim ever
        if (!hasClaimedAnyPeriod[msg.sender]) {
            hasClaimedAnyPeriod[msg.sender] = true;
            totalReferrers++;
        }

        emit RewardsClaimed(period, msg.sender, kawaiAmount);
    }

    /**
     * @notice Batch claim for multiple periods
     * @param periods Array of period numbers
     * @param kawaiAmounts Array of KAWAI amounts
     * @param merkleProofs Array of Merkle proofs
     */
    function claimMultiplePeriods(
        uint256[] calldata periods,
        uint256[] calldata kawaiAmounts,
        bytes32[][] calldata merkleProofs
    ) external nonReentrant whenNotPaused {
        require(
            periods.length == kawaiAmounts.length &&
            periods.length == merkleProofs.length,
            "Array length mismatch"
        );

        for (uint256 i = 0; i < periods.length; i++) {
            uint256 period = periods[i];
            
            if (hasClaimed[period][msg.sender]) {
                continue; // Skip already claimed periods
            }

            require(period <= currentPeriod, "Invalid period");

            // Verify Merkle proof using period-specific root
            bytes32 leaf = keccak256(
                abi.encodePacked(period, msg.sender, kawaiAmounts[i])
            );
            bytes32 periodRoot = periodMerkleRoots[period];
            require(periodRoot != bytes32(0), "Period not settled");
            require(
                MerkleProof.verify(merkleProofs[i], periodRoot, leaf),
                "Invalid proof"
            );

            // Mark as claimed
            hasClaimed[period][msg.sender] = true;

            // Mint KAWAI tokens
            IMintableToken(address(kawaiToken)).mint(msg.sender, kawaiAmounts[i]);
            totalKawaiDistributed += kawaiAmounts[i];

            emit RewardsClaimed(period, msg.sender, kawaiAmounts[i]);
        }
    }

    /**
     * @notice Update Merkle root for current period
     * @param _merkleRoot New Merkle root
     */
    function setMerkleRoot(bytes32 _merkleRoot) external onlyOwner {
        emit MerkleRootUpdated(currentPeriod, merkleRoot, _merkleRoot);
        merkleRoot = _merkleRoot;
        periodMerkleRoots[currentPeriod] = _merkleRoot;
    }
    
    /**
     * @notice Set Merkle root for a specific period (for retroactive settlements)
     * @param period Period number
     * @param _merkleRoot Merkle root for the period
     */
    function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) external onlyOwner {
        require(period <= currentPeriod, "Cannot set future period");
        emit MerkleRootUpdated(period, periodMerkleRoots[period], _merkleRoot);
        periodMerkleRoots[period] = _merkleRoot;
    }

    /**
     * @notice Advance to next period (weekly/monthly)
     * @param _merkleRoot Merkle root for new period
     */
    function advancePeriod(bytes32 _merkleRoot) external onlyOwner {
        uint256 oldPeriod = currentPeriod;
        currentPeriod++;
        merkleRoot = _merkleRoot;
        periodMerkleRoots[currentPeriod] = _merkleRoot;
        
        emit PeriodAdvanced(oldPeriod, currentPeriod);
        emit MerkleRootUpdated(currentPeriod, bytes32(0), _merkleRoot);
    }


    /**
     * @notice Check if user has claimed for a specific period
     * @param period Period number
     * @param user User address
     * @return bool Whether user has claimed
     */
    function hasClaimedPeriod(
        uint256 period,
        address user
    ) external view returns (bool) {
        return hasClaimed[period][user];
    }

    /**
     * @notice Get contract statistics
     * @return period Current period
     * @return kawaiDistributed Total KAWAI distributed
     * @return referrers Total unique referrers
     */
    function getStats() external view returns (
        uint256 period,
        uint256 kawaiDistributed,
        uint256 referrers
    ) {
        return (
            currentPeriod,
            totalKawaiDistributed,
            totalReferrers
        );
    }
    
    // ============ Emergency Functions ============
    
    /**
     * @notice Pause all claim operations (emergency only)
     * @dev Only owner can pause. Use when critical bug found or security incident.
     */
    function pause() external onlyOwner {
        _pause();
    }
    
    /**
     * @notice Unpause claim operations
     * @dev Only owner can unpause. Use after issue is resolved.
     */
    function unpause() external onlyOwner {
        _unpause();
    }
}

