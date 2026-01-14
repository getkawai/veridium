// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

interface IMintableToken {
    function mint(address to, uint256 amount) external;
}

/**
 * @title DepositCashbackDistributor
 * @dev Distributes KAWAI cashback rewards for USDT deposits using Merkle proofs
 * 
 * Flow:
 * 1. User deposits USDT → Backend calculates cashback (off-chain)
 * 2. Weekly: Backend generates Merkle tree for all pending cashback
 * 3. Owner sets Merkle root for the period
 * 4. Users claim their cashback using Merkle proof
 * 
 * Features:
 * - Period-based claims (weekly settlements)
 * - Multi-period batch claims
 * - Gas-efficient (users only pay gas when claiming)
 * - Consistent with contributor/referral reward systems
 * - Emergency pause mechanism for security
 */
contract DepositCashbackDistributor is Ownable, ReentrancyGuard, Pausable {
    
    // ============ State Variables ============
    
    IERC20 public immutable kawaiToken;
    
    // Current period (incremented weekly)
    uint256 public currentPeriod;
    
    // Global Merkle root (for backward compatibility)
    bytes32 public merkleRoot;
    
    // Period-specific Merkle roots
    mapping(uint256 => bytes32) public periodMerkleRoots;
    
    // Claim tracking: period => user => claimed
    mapping(uint256 => mapping(address => bool)) public hasClaimed;
    
    // Track unique users (for stats)
    mapping(address => bool) public hasClaimedAnyPeriod;
    
    // Statistics
    uint256 public totalKawaiDistributed;
    uint256 public totalUsers;
    
    // Allocation tracking
    uint256 public constant TOTAL_ALLOCATION = 200_000_000 * 1e18; // 200M KAWAI
    
    // ============ Events ============
    
    event CashbackClaimed(
        uint256 indexed period,
        address indexed user,
        uint256 kawaiAmount
    );
    
    event MerkleRootUpdated(
        uint256 indexed period,
        bytes32 oldRoot,
        bytes32 newRoot
    );
    
    event PeriodAdvanced(
        uint256 indexed newPeriod,
        bytes32 merkleRoot
    );
    
    // ============ Constructor ============
    
    constructor(address kawaiToken_) Ownable(msg.sender) {
        require(kawaiToken_ != address(0), "Invalid KAWAI address");
        kawaiToken = IERC20(kawaiToken_);
        currentPeriod = 1;
    }
    
    // ============ Main Functions ============
    
    /**
     * @notice Claim cashback for a specific period
     * @param period Period number (1, 2, 3, ...)
     * @param kawaiAmount Amount of KAWAI to claim
     * @param merkleProof Merkle proof for verification
     */
    function claimCashback(
        uint256 period,
        uint256 kawaiAmount,
        bytes32[] calldata merkleProof
    ) external nonReentrant whenNotPaused {
        require(period <= currentPeriod, "Invalid period");
        require(!hasClaimed[period][msg.sender], "Already claimed for this period");
        require(kawaiAmount > 0, "No cashback to claim");
        require(
            totalKawaiDistributed + kawaiAmount <= TOTAL_ALLOCATION,
            "Exceeds total allocation"
        );
        
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
        
        // Track unique users
        if (!hasClaimedAnyPeriod[msg.sender]) {
            hasClaimedAnyPeriod[msg.sender] = true;
            totalUsers++;
        }
        
        emit CashbackClaimed(period, msg.sender, kawaiAmount);
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
        
        uint256 totalAmount = 0;
        
        for (uint256 i = 0; i < periods.length; i++) {
            uint256 period = periods[i];
            uint256 amount = kawaiAmounts[i];
            
            // Skip if already claimed
            if (hasClaimed[period][msg.sender]) {
                continue;
            }
            
            require(period <= currentPeriod, "Invalid period");
            require(amount > 0, "No cashback to claim");
            
            // Verify Merkle proof
            bytes32 leaf = keccak256(
                abi.encodePacked(period, msg.sender, amount)
            );
            bytes32 periodRoot = periodMerkleRoots[period];
            require(periodRoot != bytes32(0), "Period not settled");
            require(
                MerkleProof.verify(merkleProofs[i], periodRoot, leaf),
                "Invalid proof"
            );
            
            // Mark as claimed
            hasClaimed[period][msg.sender] = true;
            totalAmount += amount;
            
            emit CashbackClaimed(period, msg.sender, amount);
        }
        
        require(totalAmount > 0, "No cashback to claim");
        require(
            totalKawaiDistributed + totalAmount <= TOTAL_ALLOCATION,
            "Exceeds total allocation"
        );
        
        // Mint total KAWAI tokens
        IMintableToken(address(kawaiToken)).mint(msg.sender, totalAmount);
        totalKawaiDistributed += totalAmount;
        
        // Track unique users
        if (!hasClaimedAnyPeriod[msg.sender]) {
            hasClaimedAnyPeriod[msg.sender] = true;
            totalUsers++;
        }
    }
    
    // ============ View Functions ============
    
    /**
     * @notice Check if user has claimed for a specific period
     */
    function hasUserClaimed(uint256 period, address user) external view returns (bool) {
        return hasClaimed[period][user];
    }
    
    /**
     * @notice Get stats
     */
    function getStats() external view returns (
        uint256 _currentPeriod,
        uint256 _totalKawaiDistributed,
        uint256 _remainingAllocation,
        uint256 _totalUsers
    ) {
        return (
            currentPeriod,
            totalKawaiDistributed,
            TOTAL_ALLOCATION - totalKawaiDistributed,
            totalUsers
        );
    }
    
    /**
     * @notice Get Merkle root for a specific period
     */
    function getPeriodMerkleRoot(uint256 period) external view returns (bytes32) {
        return periodMerkleRoots[period];
    }
    
    // ============ Admin Functions ============
    
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
     * @notice Advance to next period
     * @param _merkleRoot Merkle root for new period
     */
    function advancePeriod(bytes32 _merkleRoot) external onlyOwner {
        currentPeriod++;
        merkleRoot = _merkleRoot;
        periodMerkleRoots[currentPeriod] = _merkleRoot;
        emit PeriodAdvanced(currentPeriod, _merkleRoot);
    }
    
    /**
     * @notice Set Merkle root for a specific past period (correction only)
     * @param period Period number
     * @param _merkleRoot Merkle root
     */
    function setPeriodMerkleRoot(uint256 period, bytes32 _merkleRoot) external onlyOwner {
        require(period <= currentPeriod, "Invalid period");
        emit MerkleRootUpdated(period, periodMerkleRoots[period], _merkleRoot);
        periodMerkleRoots[period] = _merkleRoot;
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

