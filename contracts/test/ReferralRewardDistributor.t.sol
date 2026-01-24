// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Test.sol";
import "../contracts/ReferralRewardDistributor.sol";
import "../contracts/KawaiToken.sol";

contract ReferralRewardDistributorTest is Test {
    ReferralRewardDistributor public distributor;
    KawaiToken public kawai;
    
    address public owner = address(this);
    address public user1 = address(0x1);
    address public user2 = address(0x2);
    address public user3 = address(0x3);
    
    // Test amounts
    uint256 constant KAWAI_REWARD = 100 * 10**18; // 100 KAWAI
    
    function setUp() public {
        // Deploy tokens
        kawai = new KawaiToken(owner, owner);
        
        // Deploy distributor (KAWAI only, no USDT)
        distributor = new ReferralRewardDistributor(
            address(kawai)
        );
        
        // Grant MINTER_ROLE to distributor
        kawai.grantRole(kawai.MINTER_ROLE(), address(distributor));
    }
    
    function testClaimRewards() public {
        // Create Merkle tree for single user
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),      // period
            user1,           // user
            KAWAI_REWARD     // kawai amount only
        ));
        
        bytes32[] memory proof = new bytes32[](0); // Empty proof for single leaf
        bytes32 root = leaf;
        
        distributor.setMerkleRoot(root);
        
        // User claims rewards
        vm.prank(user1);
        distributor.claimRewards(1, KAWAI_REWARD, proof);
        
        // Verify balances
        assertEq(kawai.balanceOf(user1), KAWAI_REWARD);
        
        // Verify claimed status
        assertTrue(distributor.hasClaimed(1, user1));
    }
    
    function testCannotClaimTwice() public {
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            user1,
            KAWAI_REWARD
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        bytes32 root = leaf;
        
        distributor.setMerkleRoot(root);
        
        // First claim succeeds
        vm.prank(user1);
        distributor.claimRewards(1, KAWAI_REWARD, proof);
        
        // Second claim fails
        vm.prank(user1);
        vm.expectRevert("Already claimed for this period");
        distributor.claimRewards(1, KAWAI_REWARD, proof);
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
    
    function testGetStats() public {
        // Initial stats
        (
            uint256 period,
            uint256 kawaiDistributed,
            uint256 referrers
        ) = distributor.getStats();
        
        assertEq(period, 1);
        assertEq(kawaiDistributed, 0);
        assertEq(referrers, 0);
        
        // Claim rewards
        bytes32 leaf = keccak256(abi.encodePacked(
            uint256(1),
            user1,
            KAWAI_REWARD
        ));
        
        bytes32[] memory proof = new bytes32[](0);
        distributor.setMerkleRoot(leaf);
        
        vm.prank(user1);
        distributor.claimRewards(1, KAWAI_REWARD, proof);
        
        // Check updated stats
        (
            ,
            kawaiDistributed,
            referrers
        ) = distributor.getStats();
        
        assertEq(kawaiDistributed, KAWAI_REWARD);
        assertEq(referrers, 1);
    }
    
    function testOnlyOwnerCanSetMerkleRoot() public {
        bytes32 newRoot = bytes32(uint256(123));
        
        vm.prank(user1);
        vm.expectRevert();
        distributor.setMerkleRoot(newRoot);
        
        // Owner can set
        distributor.setMerkleRoot(newRoot);
        assertEq(distributor.merkleRoot(), newRoot);
    }
    
    function testMultiplePeriodClaims() public {
        // Setup period 1
        bytes32 leaf1 = keccak256(abi.encodePacked(uint256(1), user1, KAWAI_REWARD));
        distributor.setMerkleRoot(leaf1);
        
        // Claim from period 1
        bytes32[] memory proof1 = new bytes32[](0);
        vm.prank(user1);
        distributor.claimRewards(1, KAWAI_REWARD, proof1);
        
        // Advance to period 2
        bytes32 leaf2 = keccak256(abi.encodePacked(uint256(2), user1, KAWAI_REWARD * 2));
        distributor.advancePeriod(leaf2);
        
        // Claim from period 2
        bytes32[] memory proof2 = new bytes32[](0);
        vm.prank(user1);
        distributor.claimRewards(2, KAWAI_REWARD * 2, proof2);
        
        // Verify total balance
        assertEq(kawai.balanceOf(user1), KAWAI_REWARD * 3);
        
        // Verify referrers count (should be 1, not 2)
        (, , uint256 referrers) = distributor.getStats();
        assertEq(referrers, 1);
    }
}
