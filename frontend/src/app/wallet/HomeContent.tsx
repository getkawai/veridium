import {
  Card,
  Modal,
  Table,
  Tag,
  Button,
  Empty,
  Popover,
  Spin,
  Tooltip,
} from "antd";
import { memo, useState } from "react";
import { Browser } from "@wailsio/runtime";
import { JarvisService } from "@@/github.com/kawai-network/veridium/internal/services";
import {
  History,
  Plus,
  Send,
  Eye,
  EyeOff,
  Gift,
  Repeat2,
  Fuel,
  ArrowDownToLine,
  Coins,
  RefreshCw,
} from "lucide-react";
import { Icon } from "@lobehub/ui";
import { Flexbox } from "react-layout-kit";
import { useTheme } from "antd-style";
import { NetworkIcon } from "./NetworkIcons";
import { StablecoinIcon } from "./StablecoinIcon";
import type { HomeContentProps, NetworkInfo } from "./types";

// TypeScript interface for transaction analysis
interface TransactionAnalysis {
  status: "done" | "reverted" | "pending";
  txType?: string;
  method?: string;
  value?: string;
  gasUsed?: string;
  gasCost?: string;
  blockNumber?: number;
  params?: Array<{ name: string; type: string; value?: string }>;
  logs?: Array<{ name?: string }>;
  error?: string;
}

// Helper function to format relative time
const formatRelativeTime = (date: string | Date): string => {
  const now = new Date();
  const past = new Date(date);

  // Validate date before processing
  if (isNaN(past.getTime())) {
    return "-";
  }

  const seconds = Math.floor((now.getTime() - past.getTime()) / 1000);

  // Handle future dates
  if (seconds < 0) {
    const futureSeconds = Math.abs(seconds);
    if (futureSeconds < 60) return "Just now";
    if (futureSeconds < 3600) return `in ${Math.floor(futureSeconds / 60)}m`;
    if (futureSeconds < 86400) return `in ${Math.floor(futureSeconds / 3600)}h`;
    if (futureSeconds < 604800)
      return `in ${Math.floor(futureSeconds / 86400)}d`;
    return `in ${past.toLocaleDateString()}`;
  }

  // Handle past dates
  if (seconds < 60) return "Just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
  return past.toLocaleDateString();
};

// Helper function for safe parseFloat
const safeParseFloat = (
  value: string | number,
  defaultValue: number = 0,
): number => {
  const parsed = typeof value === "number" ? value : parseFloat(value);
  return isNaN(parsed) ? defaultValue : parsed;
};

// Helper function to get faucet URL based on network
const getFaucetUrl = (networkId?: number): string => {
  // Default faucet URLs for different networks
  // Monad testnet chain ID is 10143
  const faucetUrls: Record<number, string> = {
    10143: "https://testnet.monad.xyz/faucet", // Monad testnet
    // Add other network faucet URLs as needed
  };
  return faucetUrls[networkId ?? 10143] ?? "https://testnet.monad.xyz/faucet";
};

// Helper function to get transaction type color
const getTxTypeColor = (txType: string): string => {
  const colorMap: Record<string, string> = {
    DEPOSIT: "success",
    WITHDRAW: "error",
    SWAP: "processing",
    TRANSFER: "blue",
  };
  return colorMap[txType] || "default";
};

// Helper function to determine transaction sign
const getTxSign = (txType: string): "+" | "-" | "" => {
  if (txType === "DEPOSIT") return "+";
  if (txType === "WITHDRAW") return "-";
  return "";
};

// Shared column definitions for transaction tables
const getTransactionColumns = (
  theme: any,
  currentNetwork: NetworkInfo | null,
) => [
    {
      title: "Type",
      dataIndex: "txType",
      key: "txType",
      width: 100,
      render: (type: string) => <Tag color={getTxTypeColor(type)}>{type}</Tag>,
    },
    {
      title: "Amount",
      dataIndex: "amount",
      key: "amount",
      render: (amount: string, record: any) => {
        // Coerce amount to string before operations
        const amtStr = amount == null ? "" : String(amount);
        const displayAmount =
          amtStr.length > 15 ? `${amtStr.substring(0, 15)}...` : amtStr;
        const sign = getTxSign(record.txType);

        return (
          <Tooltip title={amtStr}>
            <span
              style={{
                color: sign === "+" ? theme.colorSuccess : theme.colorText,
                fontWeight: 600,
              }}
            >
              {sign}
              {displayAmount} USDT
            </span>
          </Tooltip>
        );
      },
    },
    {
      title: "Date",
      dataIndex: "createdAt",
      key: "createdAt",
      render: (date: string) => {
        const parsed = new Date(date);
        const tooltip = isNaN(parsed.getTime()) ? "-" : parsed.toLocaleString();
        return <Tooltip title={tooltip}>{formatRelativeTime(date)}</Tooltip>;
      },
    },
    {
      title: "TX Hash",
      dataIndex: "txHash",
      key: "txHash",
      width: 120,
      render: (txHash: string) =>
        txHash ? (
          <TransactionLink txHash={txHash} networkId={currentNetwork?.id} />
        ) : (
          "-"
        ),
    },
  ];

// Transaction Link with analysis popup
const TransactionLink = memo<{ txHash: string; networkId?: number }>(
  ({ txHash, networkId }) => {
    const theme = useTheme();
    const [analyzing, setAnalyzing] = useState(false);
    const [analysis, setAnalysis] = useState<TransactionAnalysis | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [expandedParam, setExpandedParam] = useState<number | null>(null);

    const handleAnalyze = async () => {
      // Check if analysis already exists to avoid redundant API calls
      if (!networkId || analyzing || analysis) return;
      setAnalyzing(true);
      setError(null);
      setAnalysis(null); // Clear previous analysis on retry
      try {
        const result = await JarvisService.AnalyzeTransaction(
          txHash,
          networkId,
        );
        if (result) {
          setAnalysis(result);
        } else {
          setError("No analysis data received");
        }
      } catch (e) {
        const errorMessage =
          e instanceof Error ? e.message : "Failed to analyze transaction";
        console.error("Failed to analyze transaction", e);
        setError(errorMessage);
      } finally {
        setAnalyzing(false);
      }
    };

    const shortHash = `${txHash.substring(0, 6)}...${txHash.substring(txHash.length - 4)}`;

    return (
      <Popover
        trigger="click"
        onOpenChange={(open) => open && !analysis && !error && handleAnalyze()}
        content={
          <div style={{ width: 300, maxHeight: 400, overflowY: "auto" }}>
            {analyzing ? (
              <Flexbox align="center" justify="center" style={{ padding: 20 }}>
                <Spin size="small" />
                <span style={{ marginLeft: 8 }}>Analyzing...</span>
              </Flexbox>
            ) : error ? (
              <Flexbox vertical gap={12} style={{ padding: 16 }}>
                <Flexbox align="center" gap={8}>
                  <Icon
                    icon={{ type: "fi", icon: "fi-rr-error" }}
                    size={24}
                    style={{ color: theme.colorError }}
                  />
                  <span style={{ color: theme.colorError, fontWeight: 600 }}>
                    Analysis Failed
                  </span>
                </Flexbox>
                <div
                  style={{
                    fontSize: 12,
                    color: theme.colorTextSecondary,
                    textAlign: "center",
                  }}
                >
                  {error}
                </div>
                <Button
                  type="primary"
                  size="small"
                  icon={<RefreshCw size={14} />}
                  onClick={handleAnalyze}
                  disabled={analyzing}
                >
                  Retry Analysis
                </Button>
              </Flexbox>
            ) : analysis ? (
              <Flexbox gap={12}>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  TRANSACTION ANALYSIS
                </div>

                <Flexbox gap={8}>
                  <Flexbox horizontal justify="space-between">
                    <span
                      style={{ color: theme.colorTextSecondary, fontSize: 12 }}
                    >
                      Status
                    </span>
                    <Tag
                      color={
                        analysis.status === "done"
                          ? "success"
                          : analysis.status === "reverted"
                            ? "error"
                            : "warning"
                      }
                    >
                      {analysis.status}
                    </Tag>
                  </Flexbox>

                  <Flexbox horizontal justify="space-between">
                    <span
                      style={{ color: theme.colorTextSecondary, fontSize: 12 }}
                    >
                      Type
                    </span>
                    <span style={{ fontSize: 12, fontWeight: 600 }}>
                      {analysis.txType || "Unknown"}
                    </span>
                  </Flexbox>

                  {analysis.method && (
                    <Flexbox horizontal justify="space-between">
                      <span
                        style={{
                          color: theme.colorTextSecondary,
                          fontSize: 12,
                        }}
                      >
                        Method
                      </span>
                      <Tag color="blue" style={{ fontFamily: "monospace" }}>
                        {analysis.method}
                      </Tag>
                    </Flexbox>
                  )}

                  {analysis.value && analysis.value !== "0" && (
                    <Flexbox horizontal justify="space-between">
                      <span
                        style={{
                          color: theme.colorTextSecondary,
                          fontSize: 12,
                        }}
                      >
                        Value
                      </span>
                      <span style={{ fontSize: 12 }}>{analysis.value}</span>
                    </Flexbox>
                  )}

                  {analysis.gasUsed && (
                    <Flexbox horizontal justify="space-between">
                      <span
                        style={{
                          color: theme.colorTextSecondary,
                          fontSize: 12,
                        }}
                      >
                        Gas Used
                      </span>
                      <span style={{ fontSize: 12 }}>
                        {parseInt(analysis.gasUsed, 10).toLocaleString()}
                      </span>
                    </Flexbox>
                  )}

                  {analysis.gasCost && (
                    <Flexbox horizontal justify="space-between">
                      <span
                        style={{
                          color: theme.colorTextSecondary,
                          fontSize: 12,
                        }}
                      >
                        Gas Cost
                      </span>
                      <span style={{ fontSize: 12 }}>{analysis.gasCost}</span>
                    </Flexbox>
                  )}

                  {analysis.blockNumber != null &&
                    analysis.blockNumber >= 0 && (
                      <Flexbox horizontal justify="space-between">
                        <span
                          style={{
                            color: theme.colorTextSecondary,
                            fontSize: 12,
                          }}
                        >
                          Block
                        </span>
                        <span style={{ fontSize: 12, fontFamily: "monospace" }}>
                          #{analysis.blockNumber.toLocaleString()}
                        </span>
                      </Flexbox>
                    )}
                </Flexbox>

                {/* Decoded Parameters */}
                {analysis.params && analysis.params.length > 0 && (
                  <>
                    <div
                      style={{
                        fontSize: 11,
                        color: theme.colorTextTertiary,
                        marginTop: 8,
                      }}
                    >
                      PARAMETERS
                    </div>
                    <Flexbox gap={4}>
                      {analysis.params.map((param: any, i: number) => (
                        <div
                          key={i}
                          style={{
                            padding: "4px 8px",
                            background: theme.colorFillTertiary,
                            borderRadius: 4,
                            fontSize: 11,
                          }}
                        >
                          <span style={{ color: theme.colorTextSecondary }}>
                            {param.name}
                          </span>
                          <span style={{ color: theme.colorTextTertiary }}>
                            {" "}
                            ({param.type})
                          </span>
                          <div
                            style={{
                              fontFamily: "monospace",
                              wordBreak: "break-all",
                              marginTop: 2,
                            }}
                          >
                            {expandedParam === i ||
                              (param.value?.length ?? 0) <= 50
                              ? param.value
                              : `${param.value?.substring(0, 50)}...`}
                            {(param.value?.length ?? 0) > 50 && (
                              <button
                                type="button"
                                onClick={() =>
                                  setExpandedParam(
                                    expandedParam === i ? null : i,
                                  )
                                }
                                style={{
                                  background: "none",
                                  border: "none",
                                  color: theme.colorPrimary,
                                  cursor: "pointer",
                                  fontSize: 11,
                                  padding: 0,
                                  marginLeft: 4,
                                }}
                              >
                                {expandedParam === i
                                  ? "Show less"
                                  : "Show more"}
                              </button>
                            )}
                          </div>
                        </div>
                      ))}
                    </Flexbox>
                  </>
                )}

                {/* Event Logs */}
                {analysis.logs && analysis.logs.length > 0 && (
                  <>
                    <div
                      style={{
                        fontSize: 11,
                        color: theme.colorTextTertiary,
                        marginTop: 8,
                      }}
                    >
                      EVENTS ({analysis.logs.length})
                    </div>
                    <Flexbox gap={4}>
                      {analysis.logs.slice(0, 3).map((log: any, i: number) => (
                        <Tag key={i} color="purple">
                          {log.name || "Unknown Event"}
                        </Tag>
                      ))}
                      {analysis.logs.length > 3 && (
                        <span
                          style={{
                            fontSize: 11,
                            color: theme.colorTextTertiary,
                          }}
                        >
                          +{analysis.logs.length - 3} more
                        </span>
                      )}
                    </Flexbox>
                  </>
                )}

                {analysis.error && (
                  <div
                    style={{
                      color: theme.colorError,
                      fontSize: 12,
                      marginTop: 8,
                    }}
                  >
                    Error: {analysis.error}
                  </div>
                )}
              </Flexbox>
            ) : (
              <div
                style={{
                  padding: 16,
                  textAlign: "center",
                  color: theme.colorTextSecondary,
                }}
              >
                Click to analyze transaction
              </div>
            )}
          </div>
        }
      >
        <button
          type="button"
          style={{
            fontFamily: "monospace",
            fontSize: 11,
            cursor: "pointer",
            color: theme.colorPrimary,
            textDecoration: "underline",
            background: "none",
            border: "none",
            padding: 0,
          }}
          tabIndex={0}
          aria-haspopup="dialog"
          onKeyDown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              e.stopPropagation();
            }
          }}
        >
          {shortHash}
        </button>
      </Popover>
    );
  },
);

const HomeContent = ({
  address,
  balance,
  nativeBalance,
  kawaiBalance,
  nativePrice,
  kawaiPrice,
  balanceVisible,
  setBalanceVisible,
  setModalType,
  transactions,
  styles,
  theme,
  currentNetwork,
  gasEstimate,
  currentBlock,
  balancesLoading,
}: HomeContentProps) => {
  const [showAllTx, setShowAllTx] = useState(false);

  // Calculate total portfolio value with safe parsing
  const usdtValue = safeParseFloat(balance, 0);
  const nativeValue = safeParseFloat(nativeBalance, 0) * nativePrice;
  const kawaiValue = safeParseFloat(kawaiBalance, 0) * kawaiPrice;
  const totalPortfolioValue = usdtValue + nativeValue + kawaiValue;

  // Helper for action button key press
  const handleActionKeyPress =
    (action: () => void) => (e: React.KeyboardEvent) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        action();
      }
    };

  return (
    <Flexbox style={{ maxWidth: 900, width: "100%" }} gap={20}>
      {/* Balance Card */}
      <Card className={styles.balanceCard}>
        <Tooltip title={balanceVisible ? "Hide balance" : "Show balance"}>
          <div
            className={styles.eyeButton}
            onClick={() => setBalanceVisible(!balanceVisible)}
            role="button"
            tabIndex={0}
            aria-label={balanceVisible ? "Hide balance" : "Show balance"}
            onKeyDown={handleActionKeyPress(() =>
              setBalanceVisible(!balanceVisible),
            )}
            style={{
              cursor: "pointer",
            }}
          >
            {balanceVisible ? <Eye size={16} /> : <EyeOff size={16} />}
          </div>
        </Tooltip>
        <Flexbox horizontal justify="space-between" align="center">
          <Flexbox style={{ flexDirection: "column" }} gap={4}>
            <span
              style={{
                fontSize: 11,
                color: theme.colorTextSecondary,
                textTransform: "uppercase",
                letterSpacing: "0.5px",
              }}
            >
              Total Portfolio Value
            </span>
            <div className={styles.statValue}>
              {balancesLoading ? (
                <Spin size="small" />
              ) : (
                <>
                  {balanceVisible
                    ? `$${totalPortfolioValue.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                    : "••••••"}
                  <span
                    style={{
                      fontSize: 16,
                      color: theme.colorTextTertiary,
                      marginLeft: 6,
                      fontWeight: 500,
                    }}
                  >
                    USD
                  </span>
                </>
              )}
            </div>
            {/* Native token balance */}
            {currentNetwork && (
              <div
                style={{
                  fontSize: 13,
                  color: theme.colorTextSecondary,
                  marginTop: 4,
                }}
              >
                {balanceVisible ? nativeBalance : "••••"}{" "}
                {currentNetwork.nativeTokenSymbol} (
                {nativePrice > 0 ? `$${nativePrice.toFixed(2)}` : "-"})
              </div>
            )}
          </Flexbox>
          {/* Network & Gas Info */}
          <Flexbox gap={8} align="flex-end">
            {gasEstimate && (
              <Tooltip
                title={`Max Tip: ${gasEstimate.maxTipGwei.toFixed(2)} Gwei`}
              >
                <Flexbox
                  horizontal
                  align="center"
                  gap={4}
                  style={{
                    padding: "4px 8px",
                    background: "rgba(255,255,255,0.1)",
                    borderRadius: 8,
                    fontSize: 11,
                  }}
                >
                  <Fuel size={12} />
                  <span>{gasEstimate.maxGasPriceGwei.toFixed(1)} Gwei</span>
                </Flexbox>
              </Tooltip>
            )}
            {currentBlock >= 0 && (
              <div style={{ fontSize: 10, color: theme.colorTextTertiary }}>
                Block #{currentBlock.toLocaleString()}
              </div>
            )}
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Quick Actions */}
      <Flexbox horizontal gap={12} style={{ marginTop: 4 }}>
        {[
          {
            label: "Deposit",
            icon: Plus,
            color: theme.colorSuccess,
            action: () => setModalType("deposit"),
          },
          {
            label: "Send",
            icon: Send,
            color: theme.colorInfo,
            action: () => setModalType("send"),
          },
          {
            label: "Receive",
            icon: ArrowDownToLine,
            color: theme.colorSuccess,
            action: () => setModalType("receive"),
          },
          {
            label: "Swap",
            icon: Repeat2,
            color: theme.colorWarning,
            action: () => setModalType("swap"),
          },
        ].map((item) => (
          <div
            key={item.label}
            className={styles.actionButton}
            onClick={item.action}
            role="button"
            tabIndex={0}
            aria-label={item.label}
            onKeyDown={handleActionKeyPress(item.action)}
            style={{
              cursor: "pointer",
              transition: "transform 0.2s ease, opacity 0.2s ease",
              outline: "none",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.transform = "scale(1.05)";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.transform = "scale(1)";
            }}
            onMouseDown={(e) => {
              e.currentTarget.style.transform = "scale(0.95)";
            }}
            onMouseUp={(e) => {
              e.currentTarget.style.transform = "scale(1.05)";
            }}
            onFocus={(e) => {
              e.currentTarget.style.boxShadow = `0 0 0 2px ${theme.colorPrimary}`;
            }}
            onBlur={(e) => {
              e.currentTarget.style.boxShadow = "none";
            }}
          >
            <div
              className={styles.actionCircle}
              style={{ background: `${item.color}20`, color: item.color }}
            >
              <item.icon size={24} />
            </div>
            <span style={{ fontWeight: 600, fontSize: 13 }}>{item.label}</span>
          </div>
        ))}
      </Flexbox>

      {/* Token List */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Coins size={16} /> Tokens
          </Flexbox>
        }
        size="small"
        extra={
          <Button
            type="text"
            icon={<Plus size={14} />}
            size="small"
            onClick={() => setModalType("addToken")}
          >
            Add Token
          </Button>
        }
      >
        <Flexbox gap={8}>
          {/* Native Token */}
          {currentNetwork && (
            <div
              className={styles.tokenRow}
              style={{
                cursor: "default",
                transition: "background-color 0.2s ease",
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor =
                  theme.colorFillQuaternary;
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = "transparent";
              }}
            >
              <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
                {currentNetwork && (
                  <NetworkIcon
                    name={currentNetwork.icon || "ethereum"}
                    size={24}
                    variant="mono"
                  />
                )}
                <div>
                  <div style={{ fontWeight: 600 }}>
                    {currentNetwork.nativeTokenSymbol}
                  </div>
                  <div
                    style={{ fontSize: 12, color: theme.colorTextSecondary }}
                  >
                    Native Token
                  </div>
                </div>
              </Flexbox>
              <Flexbox horizontal align="center" gap={16}>
                <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  {nativePrice > 0 ? `$${nativePrice.toFixed(2)}` : "-"}
                </span>
                <div style={{ textAlign: "right", minWidth: 70 }}>
                  <div style={{ fontWeight: 700, fontSize: 14, color: theme.colorText }}>
                    {balanceVisible ? nativeBalance : "••••"}
                  </div>
                  <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                    {balanceVisible && nativeValue > 0
                      ? `$${nativeValue.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                      : "••••"}
                  </div>
                </div>
              </Flexbox>
            </div>
          )}

          {/* Stablecoin (USDT/USDC) */}
          <div
            className={styles.tokenRow}
            style={{
              cursor: "default",
              transition: "background-color 0.2s ease",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = theme.colorFillQuaternary;
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = "transparent";
            }}
          >
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <StablecoinIcon currentNetwork={currentNetwork} size={36} />
              <div>
                <div style={{ fontWeight: 600 }}>
                  {currentNetwork?.stablecoinSymbol || 'USDT'}
                </div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  {currentNetwork?.stablecoinName || 'Tether USD'}
                </div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                $1.00
              </span>
              <div style={{ textAlign: "right", minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: theme.colorText }}>
                  {balanceVisible ? balance : "••••"}
                </div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  ${balanceVisible ? balance : "••••"}
                </div>
              </div>
            </Flexbox>
          </div>

          {/* KAWAI */}
          <div
            className={styles.tokenRow}
            style={{
              cursor: "default",
              transition: "background-color 0.2s ease",
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = theme.colorFillQuaternary;
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = "transparent";
            }}
          >
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <div
                style={{
                  width: 36,
                  height: 36,
                  borderRadius: "50%",
                  background:
                    "linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#fff",
                  fontWeight: 800,
                  fontSize: 14,
                  fontFamily: "Arial, sans-serif",
                  boxShadow: "0 2px 8px rgba(255, 154, 158, 0.3)",
                }}
              >
                <Icon icon={Gift} size={20} color="#fff" />
              </div>
              <div>
                <div style={{ fontWeight: 600 }}>KAWAI</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  Kawai Token
                </div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                {kawaiPrice > 0 ? `$${kawaiPrice.toFixed(4)}` : "-"}
              </span>
              <div style={{ textAlign: "right", minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: theme.colorText }}>
                  {balanceVisible ? kawaiBalance : "••••"}
                </div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  {balanceVisible && kawaiValue > 0
                    ? `$${kawaiValue.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
                    : "••••"}
                </div>
              </div>
            </Flexbox>
          </div>
        </Flexbox>
      </Card>

      {/* Activity */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <History size={16} /> Recent Activity
          </Flexbox>
        }
        size="small"
        extra={
          transactions.length > 5 && (
            <Button type="link" size="small" onClick={() => setShowAllTx(true)}>
              View All ({transactions.length})
            </Button>
          )
        }
      >
        {transactions.length > 0 ? (
          <Table
            dataSource={transactions.slice(0, 5)}
            rowKey="id"
            pagination={false}
            size="small"
            columns={getTransactionColumns(theme, currentNetwork)}
          />
        ) : (
          <Flexbox
            vertical
            align="center"
            gap={16}
            style={{ padding: "24px 0" }}
          >
            <Empty
              description={
                <span style={{ color: theme.colorTextSecondary }}>
                  No transactions yet
                </span>
              }
              image={Empty.PRESENTED_IMAGE_SIMPLE}
            />
            {currentNetwork?.isTestnet && (
              <Flexbox vertical align="center" gap={8}>
                <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>
                  Get test tokens to start exploring
                </span>
                <Button
                  type="primary"
                  size="small"
                  icon={<ArrowDownToLine size={14} />}
                  onClick={() =>
                    Browser.OpenURL(getFaucetUrl(currentNetwork?.id))
                  }
                >
                  Get Test Tokens (Faucet)
                </Button>
              </Flexbox>
            )}
          </Flexbox>
        )}
      </Card>

      {/* Full Transaction History Modal */}
      <Modal
        title={
          <Flexbox horizontal align="center" gap={8}>
            <History size={18} /> Transaction History
          </Flexbox>
        }
        open={showAllTx}
        onCancel={() => setShowAllTx(false)}
        footer={null}
        width={700}
        styles={{
          body: { padding: 16 },
        }}
      >
        <Table
          dataSource={transactions}
          rowKey="id"
          pagination={{
            pageSize: 10,
            showSizeChanger: false,
            showTotal: (total) => `${total} transactions`,
          }}
          size="small"
          columns={getTransactionColumns(theme, currentNetwork)}
        />
      </Modal>
    </Flexbox>
  );
};

export default HomeContent;
