// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Test.sol";
import "../contracts/MiningRewardDistributor.sol";
import "../contracts/KawaiToken.sol";

contract MiningRewardDistributorTest is Test {
    MiningRewardDistributor public distributor;
    KawaiToken public kawai;
    
    address public owner = address(this);
    address public developer = address(0xdE7);
    address public contributor = address(0x1);
    address public user = address(0x2);
    address public affiliator = address(0x3);
    
    // Event declaration (must match contract)
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
    
    // Test amounts (in wei, 18 decimals)
    uint256 constant BASE_REWARD = 1000 * 10**18; // 1000 KAWAI base reward
    
    // Referral user split (85/5/5/5)
    uint256 constant CONTRIBUTOR_REFERRAL = 850 * 10**18;  // 85%
    uint256 constant DEVELOPER_SHARE = 50 * 10**18;        // 5%
    uint256 constant USER_CASHBACK = 50 * 10**18;          // 5%
    uint256 constant AFFILIATOR_COMMISSION = 50 * 10**18;  // 5%
    
    // Non-referral user split (90/5/5)
    uint256 constant CONTRIBUTOR_NON_REFERRAL = 900 * 10**18; // 90%
    // DEVELOPER_SHARE and USER_CASHBACK same as above
    
    function setUp() public {
        // Deploy tokens
        kawai = new KawaiToken(owner, owner);
        
        // Deploy distributor
        distributor = new MiningRewardDistributor(address(kawai));
        
        // Grant MINTER_ROLE to distributor
        kawai.grantRole(kawai.MINTER_ROLE(), address(distributor));
    }
    
    function testClaimReferralUserReward() public {
        // Create Merkle leaf for referral user (85/5/5/5 split) - 9 fields
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),                 // period
            contributor,                // contributor
            CONTRIBUTOR_REFERRAL,       // 850 KAWAI
            DEVELOPER_SHARE,            // 50 KAWAI
            USER_CASHBACK,              // 50 KAWAI
            AFFILIATOR_COMMISSION,      // 50 KAWAI
            developer,                  // developer address (from GetRandomTreasuryAddress)
            user,                       // user address
            affiliator                  // affiliator address
        ));
        
        bytes32[] memory proof = new bytes32[](0); // Empty proof for single leaf
        bytes32 root = leaf;
        
        distributor.setMerkleRoot(root);
        
        // Contributor claims (distributes to all parties)
        vm.prank(contributor);
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
        
        // Verify all balances
        assertEq(kawai.balanceOf(contributor), CONTRIBUTOR_REFERRAL, "Contributor balance incorrect");
        assertEq(kawai.balanceOf(developer), DEVELOPER_SHARE, "Developer balance incorrect");
        assertEq(kawai.balanceOf(user), USER_CASHBACK, "User cashback incorrect");
        assertEq(kawai.balanceOf(affiliator), AFFILIATOR_COMMISSION, "Affiliator commission incorrect");
        
        // Verify claimed status
        assertTrue(distributor.hasClaimed(1, contributor));
        
        // Verify stats
        (
            uint256 period,
            uint256 contributorRewards,
            uint256 developerRewards,
            uint256 userRewards,
            uint256 affiliatorRewards
        ) = distributor.getStats();
        
        assertEq(period, 1);
        assertEq(contributorRewards, CONTRIBUTOR_REFERRAL);
        assertEq(developerRewards, DEVELOPER_SHARE);
        assertEq(userRewards, USER_CASHBACK);
        assertEq(affiliatorRewards, AFFILIATOR_COMMISSION);
    }
    
    function testClaimNonReferralUserReward() public {
        // Create Merkle leaf for non-referral user (90/5/5 split, no affiliator) - 9 fields
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),                     // period
            contributor,                    // contributor
            CONTRIBUTOR_NON_REFERRAL,       // 900 KAWAI
            DEVELOPER_SHARE,                // 50 KAWAI
            USER_CASHBACK,                  // 50 KAWAI
            uint256(0),                     // 0 KAWAI for affiliator
            developer,                      // developer address
            user,                           // user address
            address(0)                      // no affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        bytes32 root = leaf;
        
        distributor.setMerkleRoot(root);
        
        // Contributor claims
        vm.prank(contributor);
        distributor.claimReward(
            1,
            CONTRIBUTOR_NON_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            0,              // no affiliator reward
            developer,
            user,
            address(0),     // no affiliator
            proof
        );
        
        // Verify balances
        assertEq(kawai.balanceOf(contributor), CONTRIBUTOR_NON_REFERRAL, "Contributor balance incorrect");
        assertEq(kawai.balanceOf(developer), DEVELOPER_SHARE, "Developer balance incorrect");
        assertEq(kawai.balanceOf(user), USER_CASHBACK, "User cashback incorrect");
        assertEq(kawai.balanceOf(address(0)), 0, "Affiliator should have 0");
        
        // Verify stats (no affiliator rewards)
        (
            ,
            uint256 contributorRewards,
            uint256 developerRewards,
            uint256 userRewards,
            uint256 affiliatorRewards
        ) = distributor.getStats();
        
        assertEq(contributorRewards, CONTRIBUTOR_NON_REFERRAL);
        assertEq(developerRewards, DEVELOPER_SHARE);
        assertEq(userRewards, USER_CASHBACK);
        assertEq(affiliatorRewards, 0, "Should have no affiliator rewards");
    }
    
    function testCannotClaimTwice() public {
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        bytes32 root = leaf;
        
        distributor.setMerkleRoot(root);
        
        // First claim succeeds
        vm.prank(contributor);
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
        
        // Second claim fails
        vm.prank(contributor);
        vm.expectRevert("Already claimed for this period");
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
    }
    
    function testInvalidProof() public {
        bytes32 correctLeaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(correctLeaf);
        
        // Try to claim with wrong amounts (invalid proof)
        vm.prank(contributor);
        vm.expectRevert("Invalid proof");
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL + 100, // wrong amount
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
    }
    
    function testBatchClaimMultiplePeriods() public {
        // Setup 3 periods
        uint256[] memory periods = new uint256[](3);
        uint256[] memory contributorAmounts = new uint256[](3);
        uint256[] memory developerAmounts = new uint256[](3);
        uint256[] memory userAmounts = new uint256[](3);
        uint256[] memory affiliatorAmounts = new uint256[](3);
        address[] memory users = new address[](3);
        address[] memory affiliators = new address[](3);
        bytes32[][] memory proofs = new bytes32[][](3);
        
        for (uint256 i = 0; i < 3; i++) {
            periods[i] = i + 1;
            contributorAmounts[i] = CONTRIBUTOR_REFERRAL;
            developerAmounts[i] = DEVELOPER_SHARE;
            userAmounts[i] = USER_CASHBACK;
            affiliatorAmounts[i] = AFFILIATOR_COMMISSION;
            users[i] = user;
            affiliators[i] = affiliator;
            proofs[i] = new bytes32[](0);
            
            // Create leaf and set root for each period (9 fields)
            bytes32 leaf = keccak256(abi.encodePacked(
                periods[i],
                contributor,
                contributorAmounts[i],
                developerAmounts[i],
                userAmounts[i],
                affiliatorAmounts[i],
                developer,
                users[i],
                affiliators[i]
            ));
            
            if (i == 0) {
                distributor.setMerkleRoot(leaf);
            } else {
                distributor.advancePeriod(leaf);
            }
        }
        
        // Create developers array
        address[] memory developers = new address[](3);
        for (uint256 i = 0; i < 3; i++) {
            developers[i] = developer;
        }
        
        // Batch claim all 3 periods
        vm.prank(contributor);
        distributor.claimMultiplePeriods(
            periods,
            contributorAmounts,
            developerAmounts,
            userAmounts,
            affiliatorAmounts,
            developers,
            users,
            affiliators,
            proofs
        );
        
        // Verify total balances (3x each amount)
        assertEq(kawai.balanceOf(contributor), CONTRIBUTOR_REFERRAL * 3);
        assertEq(kawai.balanceOf(developer), DEVELOPER_SHARE * 3);
        assertEq(kawai.balanceOf(user), USER_CASHBACK * 3);
        assertEq(kawai.balanceOf(affiliator), AFFILIATOR_COMMISSION * 3);
        
        // Verify all periods claimed
        assertTrue(distributor.hasClaimed(1, contributor));
        assertTrue(distributor.hasClaimed(2, contributor));
        assertTrue(distributor.hasClaimed(3, contributor));
    }
    
    function testAdvancePeriod() public {
        bytes32 oldRoot = bytes32(uint256(1));
        bytes32 newRoot = bytes32(uint256(2));
        
        distributor.setMerkleRoot(oldRoot);
        assertEq(distributor.currentPeriod(), 1);
        
        distributor.advancePeriod(newRoot);
        assertEq(distributor.currentPeriod(), 2);
        assertEq(distributor.merkleRoot(), newRoot);
    }
    
    function testOnlyOwnerCanSetMerkleRoot() public {
        bytes32 root = bytes32(uint256(1));
        
        vm.prank(user);
        vm.expectRevert();
        distributor.setMerkleRoot(root);
    }
    
    function testOnlyOwnerCanAdvancePeriod() public {
        bytes32 root = bytes32(uint256(1));
        
        vm.prank(user);
        vm.expectRevert();
        distributor.advancePeriod(root);
    }
    
    function testInvalidUserAddress() public {
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            address(0),  // invalid user address
            affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(leaf);
        
        vm.prank(contributor);
        vm.expectRevert("Invalid user address");
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            address(0),  // invalid user
            affiliator,
            proof
        );
    }
    
    function testZeroAmountsHandled() public {
        // Test with all zero amounts (edge case)
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            uint256(0),
            uint256(0),
            uint256(0),
            uint256(0),
            developer,
            user,
            address(0)
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(leaf);
        
        vm.prank(contributor);
        distributor.claimReward(
            1,
            0,
            0,
            0,
            0,
            developer,
            user,
            address(0),
            proof
        );
        
        // Should succeed but no tokens minted
        assertEq(kawai.balanceOf(contributor), 0);
        assertEq(kawai.balanceOf(developer), 0);
        assertEq(kawai.balanceOf(user), 0);
        
        // But should still mark as claimed
        assertTrue(distributor.hasClaimed(1, contributor));
    }
    
    function testEventEmission() public {
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(leaf);
        
        // Expect RewardClaimed event (9 parameters)
        vm.expectEmit(true, true, true, true);
        emit RewardClaimed(
            1,
            contributor,
            user,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            affiliator
        );
        
        vm.prank(contributor);
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
    }
    
    function testGetStats() public {
        // Initial stats
        (
            uint256 period,
            uint256 contributorRewards,
            uint256 developerRewards,
            uint256 userRewards,
            uint256 affiliatorRewards
        ) = distributor.getStats();
        
        assertEq(period, 1);
        assertEq(contributorRewards, 0);
        assertEq(developerRewards, 0);
        assertEq(userRewards, 0);
        assertEq(affiliatorRewards, 0);
        
        // Claim a reward
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            contributor,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(leaf);
        
        vm.prank(contributor);
        distributor.claimReward(
            1,
            CONTRIBUTOR_REFERRAL,
            DEVELOPER_SHARE,
            USER_CASHBACK,
            AFFILIATOR_COMMISSION,
            developer,
            user,
            affiliator,
            proof
        );
        
        // Check updated stats
        (
            period,
            contributorRewards,
            developerRewards,
            userRewards,
            affiliatorRewards
        ) = distributor.getStats();
        
        assertEq(period, 1);
        assertEq(contributorRewards, CONTRIBUTOR_REFERRAL);
        assertEq(developerRewards, DEVELOPER_SHARE);
        assertEq(userRewards, USER_CASHBACK);
        assertEq(affiliatorRewards, AFFILIATOR_COMMISSION);
    }
}

