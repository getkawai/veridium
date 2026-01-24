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
 * @title MiningRewardDistributor
 * @notice Distributes mining rewards with referral-based splits
 * @dev Uses Merkle proofs for gas-efficient batch distributions
 * 
 * Website: https://getkawai.com
 * Docs: https://getkawai.com/docs
 * 
 * Reward Distribution:
 * - Referral users: 85% contributor, 5% developer, 5% user, 5% affiliator
 * - Non-referral users: 90% contributor, 5% developer, 5% user
 * 
 * Features:
 * - Single claim distributes to all parties (gas efficient)
 * - Weekly Merkle settlement
 * - Batch claiming for multiple periods
 * - Full on-chain transparency via events
 * - Emergency pause mechanism for security
 */
contract MiningRewardDistributor is Ownable, ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;

    IERC20 public immutable kawaiToken;
    
    bytes32 public merkleRoot;
    uint256 public currentPeriod;
    
    // Track claimed rewards per period per contributor
    // period => contributor => claimed
    mapping(uint256 => mapping(address => bool)) public hasClaimed;
    
    // Store merkle roots per period for historical claims
    mapping(uint256 => bytes32) public periodMerkleRoots;
    
    // Track total rewards distributed
    uint256 public totalContributorRewards;
    uint256 public totalDeveloperRewards;
    uint256 public totalUserRewards;
    uint256 public totalAffiliatorRewards;

    event RewardClaimed(
        uint256 indexed period,
        address indexed contributor,
        address indexed user,
        uint256 contributorAmount,
        uint256 developerAmount,
        uint256 userAmount,
        uint256 affiliatorAmount,
        address developer,
        address affiliator
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
    constructor(address kawaiToken_) Ownable(msg.sender) {
        require(kawaiToken_ != address(0), "Invalid KAWAI address");
        
        kawaiToken = IERC20(kawaiToken_);
        currentPeriod = 1;
    }

    /**
     * @notice Claim mining rewards with referral split
     * @dev Contributor initiates claim, all parties receive rewards in one transaction
     * @param period The settlement period
     * @param contributorAmount Amount for contributor (85% or 90%)
     * @param developerAmount Amount for developer (5%)
     * @param userAmount Amount for user/requester (5%)
     * @param affiliatorAmount Amount for affiliator (5% or 0%)
     * @param developer Developer/treasury address (from GetRandomTreasuryAddress)
     * @param user User/requester address
     * @param affiliator Referrer address (zero address if non-referral)
     * @param merkleProof Merkle proof for verification
     */
    function claimReward(
        uint256 period,
        uint256 contributorAmount,
        uint256 developerAmount,
        uint256 userAmount,
        uint256 affiliatorAmount,
        address developer,
        address user,
        address affiliator,
        bytes32[] calldata merkleProof
    ) external nonReentrant whenNotPaused {
        require(period <= currentPeriod, "Invalid period");
        require(!hasClaimed[period][msg.sender], "Already claimed for this period");
        require(user != address(0), "Invalid user address");
        
        // Verify Merkle proof using period-specific root (9-field leaf)
        bytes32 leaf = keccak256(
            abi.encodePacked(
                period,
                msg.sender,          // contributor
                contributorAmount,
                developerAmount,
                userAmount,
                affiliatorAmount,
                developer,           // developer address (from GetRandomTreasuryAddress)
                user,
                affiliator
            )
        );
        
        bytes32 periodRoot = periodMerkleRoots[period];
        require(periodRoot != bytes32(0), "Period not settled");
        require(
            MerkleProof.verify(merkleProof, periodRoot, leaf),
            "Invalid proof"
        );

        // Mark as claimed
        hasClaimed[period][msg.sender] = true;

        // Distribute rewards to all parties
        if (contributorAmount > 0) {
            IMintableToken(address(kawaiToken)).mint(msg.sender, contributorAmount);
            totalContributorRewards += contributorAmount;
        }
        
        if (developerAmount > 0 && developer != address(0)) {
            IMintableToken(address(kawaiToken)).mint(developer, developerAmount);
            totalDeveloperRewards += developerAmount;
        }
        
        if (userAmount > 0) {
            IMintableToken(address(kawaiToken)).mint(user, userAmount);
            totalUserRewards += userAmount;
        }
        
        if (affiliatorAmount > 0 && affiliator != address(0)) {
            IMintableToken(address(kawaiToken)).mint(affiliator, affiliatorAmount);
            totalAffiliatorRewards += affiliatorAmount;
        }

        emit RewardClaimed(
            period,
            msg.sender,
            user,
            contributorAmount,
            developerAmount,
            userAmount,
            affiliatorAmount,
            developer,
            affiliator
        );
    }

    /**
     * @notice Batch claim for multiple periods with silent-skip pattern
     * @dev Already-claimed periods are silently skipped without reverting.
     * This allows partial success in batch claims, useful for retry scenarios.
     * 
     * Trade-off: Cannot distinguish "no rewards" from "already claimed" in tx result.
     * Use hasClaimed(period, contributor) to check claim status before/after.
     * 
     * @param periods Array of period numbers
     * @param contributorAmounts Array of contributor amounts
     * @param developerAmounts Array of developer amounts
     * @param userAmounts Array of user amounts
     * @param affiliatorAmounts Array of affiliator amounts
     * @param users Array of user addresses
     * @param affiliators Array of affiliator addresses
     * @param merkleProofs Array of Merkle proofs
     */
    function claimMultiplePeriods(
        uint256[] calldata periods,
        uint256[] calldata contributorAmounts,
        uint256[] calldata developerAmounts,
        uint256[] calldata userAmounts,
        uint256[] calldata affiliatorAmounts,
        address[] calldata developers,
        address[] calldata users,
        address[] calldata affiliators,
        bytes32[][] calldata merkleProofs
    ) external nonReentrant whenNotPaused {
        uint256 length = periods.length;
        require(
            length == contributorAmounts.length &&
            length == developerAmounts.length &&
            length == userAmounts.length &&
            length == affiliatorAmounts.length &&
            length == developers.length &&
            length == users.length &&
            length == affiliators.length &&
            length == merkleProofs.length,
            "Array length mismatch"
        );

        for (uint256 i = 0; i < length; ) {
            _claimSinglePeriod(
                periods[i],
                contributorAmounts[i],
                developerAmounts[i],
                userAmounts[i],
                affiliatorAmounts[i],
                developers[i],
                users[i],
                affiliators[i],
                merkleProofs[i]
            );
            
            unchecked { ++i; }
        }
    }
    
    /**
     * @dev Internal function to claim a single period (reduces stack depth)
     */
    function _claimSinglePeriod(
        uint256 period,
        uint256 contributorAmount,
        uint256 developerAmount,
        uint256 userAmount,
        uint256 affiliatorAmount,
        address developer,
        address user,
        address affiliator,
        bytes32[] calldata merkleProof
    ) private {
        // Skip already claimed periods
        if (hasClaimed[period][msg.sender]) {
            return;
        }

        require(period <= currentPeriod, "Invalid period");
        require(user != address(0), "Invalid user address");

        // Verify Merkle proof (9-field leaf)
        bytes32 leaf = keccak256(
            abi.encodePacked(
                period,
                msg.sender,
                contributorAmount,
                developerAmount,
                userAmount,
                affiliatorAmount,
                developer,
                user,
                affiliator
            )
        );
        
        bytes32 periodRoot = periodMerkleRoots[period];
        require(periodRoot != bytes32(0), "Period not settled");
        require(
            MerkleProof.verify(merkleProof, periodRoot, leaf),
            "Invalid proof"
        );

        // Mark as claimed
        hasClaimed[period][msg.sender] = true;

        // Distribute rewards
        if (contributorAmount > 0) {
            IMintableToken(address(kawaiToken)).mint(msg.sender, contributorAmount);
            totalContributorRewards += contributorAmount;
        }
        
        if (developerAmount > 0 && developer != address(0)) {
            IMintableToken(address(kawaiToken)).mint(developer, developerAmount);
            totalDeveloperRewards += developerAmount;
        }
        
        if (userAmount > 0) {
            IMintableToken(address(kawaiToken)).mint(user, userAmount);
            totalUserRewards += userAmount;
        }
        
        if (affiliatorAmount > 0 && affiliator != address(0)) {
            IMintableToken(address(kawaiToken)).mint(affiliator, affiliatorAmount);
            totalAffiliatorRewards += affiliatorAmount;
        }

        emit RewardClaimed(
            period,
            msg.sender,
            user,
            contributorAmount,
            developerAmount,
            userAmount,
            affiliatorAmount,
            developer,
            affiliator
        );
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
     * @notice Set Merkle root for a specific period (timestamp-based)
     * @param _period Period ID (timestamp)
     * @param _merkleRoot Merkle root for the period
     */
    function setMerkleRootForPeriod(uint256 _period, bytes32 _merkleRoot) external onlyOwner {
        require(_period > 0, "Invalid period");
        require(_merkleRoot != bytes32(0), "Invalid merkle root");
        
        bytes32 oldRoot = periodMerkleRoots[_period];
        periodMerkleRoots[_period] = _merkleRoot;
        
        emit MerkleRootUpdated(_period, oldRoot, _merkleRoot);
    }

    /**
     * @notice Advance to next period (weekly)
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
     * @notice Check if contributor has claimed for a specific period
     * @param period Period number
     * @param contributor Contributor address
     * @return bool Whether contributor has claimed
     */
    function hasClaimedPeriod(
        uint256 period,
        address contributor
    ) external view returns (bool) {
        return hasClaimed[period][contributor];
    }

    /**
     * @notice Get contract statistics
     * @return period Current period
     * @return contributorRewards Total contributor rewards distributed
     * @return developerRewards Total developer rewards distributed
     * @return userRewards Total user cashback distributed
     * @return affiliatorRewards Total affiliator commissions distributed
     */
    function getStats() external view returns (
        uint256 period,
        uint256 contributorRewards,
        uint256 developerRewards,
        uint256 userRewards,
        uint256 affiliatorRewards
    ) {
        return (
            currentPeriod,
            totalContributorRewards,
            totalDeveloperRewards,
            totalUserRewards,
            totalAffiliatorRewards
        );
    }

    /**
     * @notice Emergency pause - stops all claiming operations
     * @dev Only owner can pause. Use in case of security issues or bugs
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @notice Unpause contract - resumes normal operations
     * @dev Only owner can unpause. Ensure issue is resolved before unpausing
     */
    function unpause() external onlyOwner {
        _unpause();
    }
}

