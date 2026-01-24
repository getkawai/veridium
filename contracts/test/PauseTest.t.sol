// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "forge-std/Test.sol";
import {MiningRewardDistributor} from "../contracts/MiningRewardDistributor.sol";
import {DepositCashbackDistributor} from "../contracts/DepositCashbackDistributor.sol";
import {ReferralRewardDistributor} from "../contracts/ReferralRewardDistributor.sol";
import {KawaiToken} from "../contracts/KawaiToken.sol";

/**
 * @title PauseTest
 * @dev Test emergency pause mechanism for all distributor contracts
 */
contract PauseTest is Test {
    KawaiToken public kawaiToken;
    MiningRewardDistributor public miningDistributor;
    DepositCashbackDistributor public cashbackDistributor;
    ReferralRewardDistributor public referralDistributor;
    
    address public owner = address(this);
    address public user1 = address(0x1);
    
    bytes32[] public merkleProof;
    
    function setUp() public {
        // Deploy KAWAI token
        kawaiToken = new KawaiToken(owner, owner);
        
        // Deploy distributors
        miningDistributor = new MiningRewardDistributor(address(kawaiToken));
        cashbackDistributor = new DepositCashbackDistributor(address(kawaiToken));
        referralDistributor = new ReferralRewardDistributor(address(kawaiToken));
        
        // Grant MINTER_ROLE to distributors
        bytes32 MINTER_ROLE = keccak256("MINTER_ROLE");
        kawaiToken.grantRole(MINTER_ROLE, address(miningDistributor));
        kawaiToken.grantRole(MINTER_ROLE, address(cashbackDistributor));
        kawaiToken.grantRole(MINTER_ROLE, address(referralDistributor));
        
        // Setup merkle roots (simplified for testing)
        bytes32 merkleRoot = keccak256(abi.encodePacked("test"));
        miningDistributor.setMerkleRoot(merkleRoot);
        cashbackDistributor.setMerkleRoot(merkleRoot);
        referralDistributor.setMerkleRoot(merkleRoot);
    }
    
    // ============ Mining Distributor Pause Tests ============
    
    function testMiningPause() public {
        // Pause contract
        miningDistributor.pause();
        
        // Try to claim - should revert
        vm.expectRevert();
        vm.prank(user1);
        miningDistributor.claimReward(
            1,
            100 ether,
            0,
            0,
            0,
            address(0),
            user1,
            address(0),
            merkleProof
        );
    }
    
    function testMiningUnpause() public {
        // Pause then unpause
        miningDistributor.pause();
        miningDistributor.unpause();
        
        // Should be able to call (will fail on proof, but not on pause)
        vm.expectRevert("Invalid proof");
        vm.prank(user1);
        miningDistributor.claimReward(
            1,
            100 ether,
            0,
            0,
            0,
            address(0),
            user1,
            address(0),
            merkleProof
        );
    }
    
    function testMiningOnlyOwnerCanPause() public {
        vm.expectRevert();
        vm.prank(user1);
        miningDistributor.pause();
    }
    
    function testMiningOnlyOwnerCanUnpause() public {
        miningDistributor.pause();
        
        vm.expectRevert();
        vm.prank(user1);
        miningDistributor.unpause();
    }
    
    function testMiningBatchClaimPaused() public {
        // Pause contract
        miningDistributor.pause();
        
        // Try batch claim - should revert
        uint256[] memory periods = new uint256[](1);
        uint256[] memory contributorAmounts = new uint256[](1);
        uint256[] memory developerAmounts = new uint256[](1);
        uint256[] memory userAmounts = new uint256[](1);
        uint256[] memory affiliatorAmounts = new uint256[](1);
        address[] memory developers = new address[](1);
        address[] memory users = new address[](1);
        address[] memory affiliators = new address[](1);
        bytes32[][] memory proofs = new bytes32[][](1);
        
        periods[0] = 1;
        contributorAmounts[0] = 100 ether;
        users[0] = user1;
        
        vm.expectRevert();
        vm.prank(user1);
        miningDistributor.claimMultiplePeriods(
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
    }
    
    // ============ Cashback Distributor Pause Tests ============
    
    function testCashbackPause() public {
        // Pause contract
        cashbackDistributor.pause();
        
        // Try to claim - should revert
        vm.expectRevert();
        vm.prank(user1);
        cashbackDistributor.claimCashback(
            1,
            100 ether,
            merkleProof
        );
    }
    
    function testCashbackUnpause() public {
        // Pause then unpause
        cashbackDistributor.pause();
        cashbackDistributor.unpause();
        
        // Should be able to call (will fail on proof, but not on pause)
        vm.expectRevert("Invalid proof");
        vm.prank(user1);
        cashbackDistributor.claimCashback(
            1,
            100 ether,
            merkleProof
        );
    }
    
    function testCashbackOnlyOwnerCanPause() public {
        vm.expectRevert();
        vm.prank(user1);
        cashbackDistributor.pause();
    }
    
    function testCashbackOnlyOwnerCanUnpause() public {
        cashbackDistributor.pause();
        
        vm.expectRevert();
        vm.prank(user1);
        cashbackDistributor.unpause();
    }
    
    function testCashbackBatchClaimPaused() public {
        // Pause contract
        cashbackDistributor.pause();
        
        // Try batch claim - should revert
        uint256[] memory periods = new uint256[](1);
        uint256[] memory amounts = new uint256[](1);
        bytes32[][] memory proofs = new bytes32[][](1);
        
        periods[0] = 1;
        amounts[0] = 100 ether;
        
        vm.expectRevert();
        vm.prank(user1);
        cashbackDistributor.claimMultiplePeriods(periods, amounts, proofs);
    }
    
    // ============ Referral Distributor Pause Tests ============
    
    function testReferralPause() public {
        // Pause contract
        referralDistributor.pause();
        
        // Try to claim - should revert
        vm.expectRevert();
        vm.prank(user1);
        referralDistributor.claimRewards(
            1,
            100 ether,
            merkleProof
        );
    }
    
    function testReferralUnpause() public {
        // Pause then unpause
        referralDistributor.pause();
        referralDistributor.unpause();
        
        // Should be able to call (will fail on proof, but not on pause)
        vm.expectRevert("Invalid proof");
        vm.prank(user1);
        referralDistributor.claimRewards(
            1,
            100 ether,
            merkleProof
        );
    }
    
    function testReferralOnlyOwnerCanPause() public {
        vm.expectRevert();
        vm.prank(user1);
        referralDistributor.pause();
    }
    
    function testReferralOnlyOwnerCanUnpause() public {
        referralDistributor.pause();
        
        vm.expectRevert();
        vm.prank(user1);
        referralDistributor.unpause();
    }
    
    function testReferralBatchClaimPaused() public {
        // Pause contract
        referralDistributor.pause();
        
        // Try batch claim - should revert
        uint256[] memory periods = new uint256[](1);
        uint256[] memory amounts = new uint256[](1);
        bytes32[][] memory proofs = new bytes32[][](1);
        
        periods[0] = 1;
        amounts[0] = 100 ether;
        
        vm.expectRevert();
        vm.prank(user1);
        referralDistributor.claimMultiplePeriods(periods, amounts, proofs);
    }
}
