import React, { useState, useEffect } from 'react';

// Simplified ABI for the Claim function
// claim(uint256 index, address account, uint256 amount, bytes32[] merkleProof)

export const ClaimReward: React.FC = () => {
    const [account, setAccount] = useState<string | null>(null);
    const [proofData, setProofData] = useState<any>(null);
    const [loading, setLoading] = useState(false);
    const [status, setStatus] = useState<string>('');

    // Hardcoded Contract Address (TODO: Replace after deployment)
    const CONTRACT_ADDRESS = "0xYOUR_MERKLE_DISTRIBUTOR_ADDRESS";

    const connectWallet = async () => {
        if (typeof window !== 'undefined' && (window as any).ethereum) {
            try {
                const accounts = await (window as any).ethereum.request({ method: 'eth_requestAccounts' });
                setAccount(accounts[0]);
                setStatus('');
            } catch (err: any) {
                setStatus('Failed to connect wallet: ' + err.message);
            }
        } else {
            setStatus('Please install MetaMask!');
        }
    };

    const checkEligibility = async () => {
        if (!account) return;
        setLoading(true);
        setStatus('Checking eligibility...');
        try {
            // Assuming local backend for dev
            const res = await fetch(`http://localhost:8080/api/v1/claim/proof?address=${account}`);
            if (res.status === 404) {
                setStatus('No rewards found for this address.');
                setProofData(null);
            } else if (!res.ok) {
                setStatus('Error contacting server.');
            } else {
                const data = await res.json();
                setProofData(data);
                setStatus('Reward Available!');
            }
        } catch (err) {
            setStatus('Network error. Is the backend running?');
        } finally {
            setLoading(false);
        }
    };

    const handleClaim = async () => {
        setStatus('Claiming... (Note: This requires "viem" or manual ABI encoding. Check console for data)');
        console.log("Claim Data:", {
            contract: CONTRACT_ADDRESS,
            index: proofData.index,
            account: account,
            amount: proofData.amount,
            proof: proofData.proof
        });
        
        // TODO: Implement actual contract call using viem/ethers
        // Since dependencies are missing, we stop here for the UI demo.
        alert(`Prepare to claim ${Number(proofData.amount) / 1e18} KAWAI.\n\nData logged to console.\n\nPlease install 'viem' to enable in-browser transactions.`);
    };

    return (
        <div style={{ padding: 20, border: '1px solid #333', borderRadius: 8, maxWidth: 400, margin: '20px auto', background: '#111', color: '#fff', fontFamily: 'sans-serif' }}>
            <h2 style={{ borderBottom: '1px solid #444', paddingBottom: 10 }}>Contributor Reward</h2>
            
            {!account ? (
                <button 
                    onClick={connectWallet}
                    style={{ background: '#1890ff', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: 4, cursor: 'pointer', width: '100%' }}
                >
                    Connect Wallet
                </button>
            ) : (
                <div>
                    <p style={{ fontSize: 12, color: '#888' }}>Connected: {account}</p>
                    
                    {!proofData ? (
                        <button 
                            onClick={checkEligibility}
                            disabled={loading}
                            style={{ background: '#52c41a', color: '#fff', border: 'none', padding: '10px 20px', borderRadius: 4, cursor: 'pointer', width: '100%' }}
                        >
                            {loading ? 'Checking...' : 'Check Eligibility'}
                        </button>
                    ) : (
                        <div style={{ marginTop: 20 }}>
                            <div style={{ marginBottom: 20, textAlign: 'center' }}>
                                <span style={{ fontSize: 32, fontWeight: 'bold', display: 'block' }}>
                                    {(Number(proofData.amount) / 1e18).toFixed(2)}
                                </span>
                                <span style={{ color: '#aaa' }}>KAWAI Tokens</span>
                            </div>
                            
                            <button 
                                onClick={handleClaim}
                                style={{ background: '#faad14', color: '#000', border: 'none', padding: '10px 20px', borderRadius: 4, cursor: 'pointer', width: '100%', fontWeight: 'bold' }}
                            >
                                CLAIM NOW
                            </button>
                            <p style={{ fontSize: 10, color: '#666', marginTop: 10, textAlign: 'center' }}>
                                Gas fee required (BNB)
                            </p>
                        </div>
                    )}
                </div>
            )}
            
            {status && <p style={{ marginTop: 10, color: status.includes('Reward') ? '#52c41a' : '#ff4d4f', fontSize: 12 }}>{status}</p>}
        </div>
    );
};
