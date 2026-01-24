// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Test.sol";
import "../contracts/DepositCashbackDistributor.sol";
import "../contracts/KawaiToken.sol";

contract DepositCashbackDistributorTest is Test {
    DepositCashbackDistributor public distributor;
    KawaiToken public kawaiToken;
    
    address public owner = address(1);
    address public user1 = address(2);
    address public user2 = address(3);
    
    // Merkle tree data
    bytes32 public merkleRoot;
    bytes32[] public proof1;
    bytes32[] public proof2;
    
    function setUp() public {
        vm.startPrank(owner);
        
        // Deploy KAWAI token (defaultAdmin, minter)
        kawaiToken = new KawaiToken(owner, owner);
        
        // Deploy distributor
        distributor = new DepositCashbackDistributor(address(kawaiToken));
        
        // Grant MINTER_ROLE to distributor
        kawaiToken.grantRole(kawaiToken.MINTER_ROLE(), address(distributor));
        
        vm.stopPrank();
        
        // Setup Merkle tree for period 1
        // user1: 1000 KAWAI, user2: 2000 KAWAI
        setupMerkleTree();
    }
    
    function setupMerkleTree() internal {
        // Leaf 1: period=1, user1, 1000 KAWAI
        bytes32 leaf1 = keccak256(abi.encodePacked(uint256(1), user1, uint256(1000 * 1e18)));
        
        // Leaf 2: period=1, user2, 2000 KAWAI
        bytes32 leaf2 = keccak256(abi.encodePacked(uint256(1), user2, uint256(2000 * 1e18)));
        
        // Sort leaves (required for OpenZeppelin MerkleProof)
        (bytes32 left, bytes32 right) = leaf1 < leaf2 ? (leaf1, leaf2) : (leaf2, leaf1);
        
        // Calculate root (must use sorted pair)
        merkleRoot = keccak256(abi.encodePacked(left, right));
        
        // Generate proofs (sibling for each leaf)
        // NOTE: Proof assignment depends on sorted order
        // If leaf1 < leaf2: user1's proof is leaf2, user2's proof is leaf1
        // If leaf2 < leaf1: user1's proof is leaf2, user2's proof is leaf1 (same)
        proof1.push(leaf2); // user1 proof = sibling (leaf2)
        proof2.push(leaf1); // user2 proof = sibling (leaf1)
    }
    
    function testInitialState() public view {
        assertEq(distributor.currentPeriod(), 1);
        assertEq(distributor.totalKawaiDistributed(), 0);
        assertEq(distributor.totalUsers(), 0);
        assertEq(address(distributor.kawaiToken()), address(kawaiToken));
    }
    
    function testSetMerkleRoot() public {
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        assertEq(distributor.merkleRoot(), merkleRoot);
        assertEq(distributor.periodMerkleRoots(1), merkleRoot);
    }
    
    function testOnlyOwnerCanSetMerkleRoot() public {
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, user1));
        distributor.setMerkleRoot(merkleRoot);
    }
    
    function testClaimCashback() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // User1 claims
        vm.prank(user1);
        distributor.claimCashback(1, 1000 * 1e18, proof1);
        
        // Verify
        assertEq(kawaiToken.balanceOf(user1), 1000 * 1e18);
        assertEq(distributor.totalKawaiDistributed(), 1000 * 1e18);
        assertEq(distributor.totalUsers(), 1);
        assertTrue(distributor.hasClaimed(1, user1));
        assertTrue(distributor.hasClaimedAnyPeriod(user1));
    }
    
    function testCannotClaimTwice() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // User1 claims
        vm.prank(user1);
        distributor.claimCashback(1, 1000 * 1e18, proof1);
        
        // Try to claim again
        vm.prank(user1);
        vm.expectRevert("Already claimed for this period");
        distributor.claimCashback(1, 1000 * 1e18, proof1);
    }
    
    function testCannotClaimWithInvalidProof() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // Try to claim with wrong proof
        vm.prank(user1);
        vm.expectRevert("Invalid proof");
        distributor.claimCashback(1, 1000 * 1e18, proof2); // Using user2's proof
    }
    
    function testCannotClaimFuturePeriod() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // Try to claim period 2 (not yet started)
        vm.prank(user1);
        vm.expectRevert("Invalid period");
        distributor.claimCashback(2, 1000 * 1e18, proof1);
    }
    
    function testAdvancePeriod() public {
        // Set Merkle root for period 1
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // Advance to period 2
        bytes32 newRoot = keccak256("period2");
        vm.prank(owner);
        distributor.advancePeriod(newRoot);
        
        // Verify
        assertEq(distributor.currentPeriod(), 2);
        assertEq(distributor.merkleRoot(), newRoot);
        assertEq(distributor.periodMerkleRoots(2), newRoot);
        assertEq(distributor.periodMerkleRoots(1), merkleRoot); // Old root preserved
    }
    
    function testMultiplePeriodClaims() public {
        // Period 1: Set root and user1 claims
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        vm.prank(user1);
        distributor.claimCashback(1, 1000 * 1e18, proof1);
        
        // Period 2: Setup new tree
        bytes32 leaf1Period2 = keccak256(abi.encodePacked(uint256(2), user1, uint256(500 * 1e18)));
        bytes32 leaf2Period2 = keccak256(abi.encodePacked(uint256(2), user2, uint256(1500 * 1e18)));
        (bytes32 left, bytes32 right) = leaf1Period2 < leaf2Period2 
            ? (leaf1Period2, leaf2Period2) 
            : (leaf2Period2, leaf1Period2);
        bytes32 merkleRootPeriod2 = keccak256(abi.encodePacked(left, right));
        
        bytes32[] memory proof1Period2 = new bytes32[](1);
        proof1Period2[0] = leaf2Period2; // user1's proof is always leaf2 (sibling)
        
        // Advance period
        vm.prank(owner);
        distributor.advancePeriod(merkleRootPeriod2);
        
        // User1 claims period 2
        vm.prank(user1);
        distributor.claimCashback(2, 500 * 1e18, proof1Period2);
        
        // Verify
        assertEq(kawaiToken.balanceOf(user1), 1500 * 1e18); // 1000 + 500
        assertEq(distributor.totalKawaiDistributed(), 1500 * 1e18);
        assertEq(distributor.totalUsers(), 1); // Still 1 unique user
    }
    
    function testBatchClaim() public {
        // Setup period 1
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // Setup period 2
        bytes32 leaf1Period2 = keccak256(abi.encodePacked(uint256(2), user1, uint256(500 * 1e18)));
        bytes32 leaf2Period2 = keccak256(abi.encodePacked(uint256(2), user2, uint256(1500 * 1e18)));
        (bytes32 left, bytes32 right) = leaf1Period2 < leaf2Period2 
            ? (leaf1Period2, leaf2Period2) 
            : (leaf2Period2, leaf1Period2);
        bytes32 merkleRootPeriod2 = keccak256(abi.encodePacked(left, right));
        
        bytes32[] memory proof1Period2 = new bytes32[](1);
        proof1Period2[0] = leaf2Period2; // user1's proof is always leaf2 (sibling)
        
        vm.prank(owner);
        distributor.advancePeriod(merkleRootPeriod2);
        
        // Batch claim
        uint256[] memory periods = new uint256[](2);
        periods[0] = 1;
        periods[1] = 2;
        
        uint256[] memory amounts = new uint256[](2);
        amounts[0] = 1000 * 1e18;
        amounts[1] = 500 * 1e18;
        
        bytes32[][] memory proofs = new bytes32[][](2);
        proofs[0] = proof1;
        proofs[1] = proof1Period2;
        
        vm.prank(user1);
        distributor.claimMultiplePeriods(periods, amounts, proofs);
        
        // Verify
        assertEq(kawaiToken.balanceOf(user1), 1500 * 1e18);
        assertEq(distributor.totalKawaiDistributed(), 1500 * 1e18);
        assertTrue(distributor.hasClaimed(1, user1));
        assertTrue(distributor.hasClaimed(2, user1));
    }
    
    function testGetStats() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // User1 claims
        vm.prank(user1);
        distributor.claimCashback(1, 1000 * 1e18, proof1);
        
        // Get stats
        (
            uint256 currentPeriod,
            uint256 totalDistributed,
            uint256 totalUsers
        ) = distributor.getStats();
        
        assertEq(currentPeriod, 1);
        assertEq(totalDistributed, 1000 * 1e18);
        assertEq(totalUsers, 1);
    }
    
    function testMultipleUsers() public {
        // Set Merkle root
        vm.prank(owner);
        distributor.setMerkleRoot(merkleRoot);
        
        // User1 claims
        vm.prank(user1);
        distributor.claimCashback(1, 1000 * 1e18, proof1);
        
        // User2 claims
        vm.prank(user2);
        distributor.claimCashback(1, 2000 * 1e18, proof2);
        
        // Verify
        assertEq(kawaiToken.balanceOf(user1), 1000 * 1e18);
        assertEq(kawaiToken.balanceOf(user2), 2000 * 1e18);
        assertEq(distributor.totalKawaiDistributed(), 3000 * 1e18);
        assertEq(distributor.totalUsers(), 2);
    }
}

