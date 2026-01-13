# Security Checklist - Production Deployment

**CRITICAL:** Review this before deploying to production.

---

## 🔐 Smart Contract Security

### Code Review
- [ ] All contracts reviewed by 2+ developers
- [ ] No hardcoded addresses (use constructor parameters)
- [ ] No hardcoded private keys
- [ ] Access control properly implemented (onlyOwner, roles)
- [ ] Reentrancy guards on all external calls
- [ ] Integer overflow protection (Solidity 0.8+)
- [ ] Input validation on all functions

### Security Audit
- [ ] **HIGHLY RECOMMENDED:** Professional audit (CertiK/OpenZeppelin/Trail of Bits)
- [ ] Or: Comprehensive internal security review
- [ ] All findings addressed
- [ ] Audit report published (transparency)

### Access Control
- [ ] Admin functions protected (onlyOwner)
- [ ] MINTER_ROLE only granted to distributors
- [ ] Multi-sig wallet for admin (Gnosis Safe recommended)
- [ ] Timelock on critical functions (24-48 hours)
- [ ] Emergency pause mechanism implemented & tested

### Testing
- [ ] 100% test coverage on critical functions
- [ ] Fuzz testing completed
- [ ] Edge cases tested (0 amount, max uint256, etc)
- [ ] Gas optimization verified
- [ ] Upgrade mechanism tested (if applicable)

---

## 🔑 Key Management

### Private Keys
- [ ] **NEVER** commit private keys to git
- [ ] **NEVER** hardcode private keys in code
- [ ] Use environment variables for keys
- [ ] Or use secure key management (AWS KMS/HashiCorp Vault)
- [ ] Hardware wallet for admin operations (Ledger/Trezor)
- [ ] Key rotation policy (every 90 days)
- [ ] Backup keys encrypted & offline

### Admin Wallet
- [ ] Generated on secure, offline machine
- [ ] Seed phrase written down (not digital)
- [ ] Seed phrase stored in safe/vault
- [ ] Test restore from seed phrase
- [ ] Multi-sig recommended (3-of-5 or 2-of-3)
- [ ] Separate hot wallet for daily operations

### Settlement Wallet
- [ ] Dedicated wallet (not admin wallet)
- [ ] Minimal funds (only for gas)
- [ ] Automated key rotation
- [ ] Monitored for suspicious activity
- [ ] Rate limiting on transactions

---

## 🌐 Backend Security

### API Security
- [ ] Rate limiting (100 req/min per IP)
- [ ] API authentication (JWT/API keys)
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (sanitize inputs)
- [ ] CSRF protection
- [ ] CORS configured (whitelist frontend domain only)

### Infrastructure
- [ ] SSH key-only access (no password)
- [ ] Firewall configured (allow only necessary ports)
- [ ] Fail2ban installed (block brute force)
- [ ] Regular security updates (weekly)
- [ ] Intrusion detection (OSSEC/Wazuh)
- [ ] DDoS protection (Cloudflare)
- [ ] VPN for admin access (optional)

### Secrets Management
- [ ] No secrets in code
- [ ] No secrets in git
- [ ] Environment variables for secrets
- [ ] Or use AWS Secrets Manager/Vault
- [ ] Secrets rotated regularly
- [ ] Access logs for secrets

### Logging & Monitoring
- [ ] All API requests logged
- [ ] Failed authentication attempts logged
- [ ] Suspicious activity alerts
- [ ] Log retention (90 days minimum)
- [ ] Logs encrypted at rest
- [ ] SIEM integration (optional)

---

## 🖥️ Frontend Security

### Code Security
- [ ] No private keys in frontend code
- [ ] No API keys in frontend code
- [ ] Content Security Policy (CSP) configured
- [ ] Subresource Integrity (SRI) for CDN
- [ ] XSS prevention (React escapes by default)
- [ ] Dependency vulnerability scan (npm audit)

### Wallet Integration
- [ ] Use established libraries (ethers.js/web3.js)
- [ ] Never request private keys from users
- [ ] Verify transaction before signing
- [ ] Display transaction details clearly
- [ ] Warn users about phishing

### HTTPS
- [ ] SSL certificate installed
- [ ] HTTPS enforced (redirect HTTP)
- [ ] HSTS header configured
- [ ] Certificate auto-renewal (Let's Encrypt)

---

## 🔍 Monitoring & Alerts

### Security Monitoring
- [ ] Failed login attempts tracked
- [ ] Unusual transaction patterns detected
- [ ] Contract balance monitoring
- [ ] Admin function calls logged
- [ ] Suspicious wallet activity alerts

### Alerts Configuration
```yaml
Critical Alerts (Immediate):
  - Multiple failed admin logins
  - Unauthorized contract function call
  - Contract balance depleted
  - Unusual transaction volume
  - Security breach detected

Warning Alerts (1 hour):
  - High failed API requests
  - Unusual user behavior
  - Low contract balance
  - Slow API response
```

### Incident Response
- [ ] Incident response plan documented
- [ ] Team roles defined
- [ ] Emergency contacts list
- [ ] Communication templates prepared
- [ ] Post-mortem process defined
- [ ] Regular drills (quarterly)

---

## 🛡️ Operational Security

### Access Control
- [ ] Principle of least privilege
- [ ] Role-based access control (RBAC)
- [ ] Regular access review (monthly)
- [ ] Offboarding checklist (revoke access)
- [ ] Audit logs for all admin actions

### Backup & Recovery
- [ ] Daily automated backups
- [ ] Backup encryption
- [ ] Offsite backup storage
- [ ] Backup restore tested (monthly)
- [ ] Disaster recovery plan
- [ ] RTO/RPO defined (4 hours/1 hour recommended)

### Compliance
- [ ] GDPR compliance (if EU users)
- [ ] Data retention policy
- [ ] Privacy policy published
- [ ] Terms of service published
- [ ] KYC/AML if required
- [ ] Legal review completed

---

## 🚨 Red Flags (Stop Deployment)

**DO NOT DEPLOY if any of these are true:**

- [ ] Private keys in git history
- [ ] No security audit (for >$100k TVL)
- [ ] No multi-sig for admin
- [ ] No emergency pause mechanism
- [ ] No monitoring/alerts configured
- [ ] No backup strategy
- [ ] No incident response plan
- [ ] Untested on testnet
- [ ] Known critical vulnerabilities
- [ ] No access control on admin functions

---

## ✅ Security Verification Commands

Run these before deployment:

```bash
# 1. Check for secrets in git
git log -p | grep -i "private.*key\|secret\|password" || echo "✅ No secrets found"

# 2. Check for hardcoded addresses
grep -r "0x[a-fA-F0-9]{40}" --include="*.go" --include="*.sol" | grep -v "test" || echo "✅ No hardcoded addresses"

# 3. Verify contract ownership
cast call <CONTRACT> "owner()(address)"

# 4. Verify MINTER_ROLE
cast call <TOKEN> "hasRole(bytes32,address)(bool)" \
  $(cast keccak "MINTER_ROLE") <DISTRIBUTOR>

# 5. Test emergency pause
cast send <CONTRACT> "pause()" --private-key <ADMIN_KEY>
cast call <CONTRACT> "paused()(bool)"  # Should return true
cast send <CONTRACT> "unpause()" --private-key <ADMIN_KEY>

# 6. Verify rate limiting
for i in {1..150}; do curl https://api.kawai.network/health; done
# Should see 429 Too Many Requests after 100 requests
```

---

## 📚 Security Resources

### Audit Firms
- CertiK: https://www.certik.com/
- OpenZeppelin: https://openzeppelin.com/security-audits/
- Trail of Bits: https://www.trailofbits.com/
- Consensys Diligence: https://consensys.net/diligence/

### Security Tools
- Slither (static analysis): https://github.com/crytic/slither
- Mythril (symbolic execution): https://github.com/ConsenSys/mythril
- Echidna (fuzzing): https://github.com/crytic/echidna
- Manticore (symbolic execution): https://github.com/trailofbits/manticore

### Best Practices
- OpenZeppelin Contracts: https://docs.openzeppelin.com/contracts/
- Consensys Smart Contract Best Practices: https://consensys.github.io/smart-contract-best-practices/
- OWASP Top 10: https://owasp.org/www-project-top-ten/

---

## 🎯 Security Maturity Levels

### Level 1: Minimum (Testnet OK, Mainnet RISKY)
- Basic access control
- Input validation
- HTTPS enabled
- Backups configured

### Level 2: Standard (Mainnet OK for <$100k TVL)
- Internal security review
- Multi-sig for admin
- Monitoring & alerts
- Incident response plan
- Regular updates

### Level 3: Advanced (Mainnet OK for >$100k TVL)
- Professional security audit
- Timelock on admin functions
- Emergency pause mechanism
- Bug bounty program
- Regular penetration testing
- 24/7 monitoring

### Level 4: Enterprise (Mainnet OK for >$1M TVL)
- Multiple security audits
- Formal verification
- Insurance coverage
- Dedicated security team
- Real-time threat detection
- Compliance certifications

---

**Current Status:** [Fill after review]

**Target Level:** Level 2 (Standard) minimum

**Gaps:** [List any gaps]

**Action Plan:** [Plan to address gaps]

---

**Remember:** Security is not a one-time task. It's an ongoing process.

**Review this checklist:** Before deployment, monthly, and after any major changes.
