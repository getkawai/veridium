'use client';

import { Alert, Space, Typography } from 'antd';
import { SafetyOutlined, LockOutlined, CloudOffOutlined } from '@ant-design/icons';
import { createStyles } from 'antd-style';
import { memo } from 'react';

const { Text } = Typography;

export const useStyles = createStyles(({ css, token }) => ({
  container: css`
    padding: 16px;
    background: linear-gradient(135deg, ${token.colorPrimaryBg} 0%, ${token.colorBgContainer} 100%);
    border-bottom: 1px solid ${token.colorBorder};
  `,
  title: css`
    margin: 0 0 8px 0 !important;
    font-size: 16px;
    font-weight: 600;
    color: ${token.colorText};
  `,
  badge: css`
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 12px;
    background: ${token.colorSuccessBg};
    color: ${token.colorSuccessText};
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  `,
  features: css`
    display: flex;
    gap: 24px;
    flex-wrap: wrap;
    margin-top: 12px;
  `,
  feature: css`
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: ${token.colorTextSecondary};
  `,
}));

const PrivacyBanner = memo(() => {
  const { styles } = useStyles();

  return (
    <div className={styles.container}>
      <Alert
        message={
          <div>
            <Space direction="vertical" style={{ width: '100%' }}>
              <div>
                <Text strong className={styles.title}>
                  🔒 "Own Your Data, Command Your AI" • Kawai Privacy Commitment
                </Text>
                <div className={styles.badge}>
                  <SafetyOutlined /> LOCAL AI • PRIVACY FIRST
                </div>
              </div>

              <div className={styles.features}>
                <div className={styles.feature}>
                  <LockOutlined />
                  <Text>100% Local - No Cloud Dependency</Text>
                </div>
                <div className={styles.feature}>
                  <CloudOffOutlined />
                  <Text>Full Data Sovereignty</Text>
                </div>
                <div className={styles.feature}>
                  <SafetyOutlined />
                  <Text>Zero Tracking, Zero Sharing</Text>
                </div>
              </div>

              <Text type="secondary" style={{ fontSize: '12px', marginTop: 8 }}>
                All database operations are performed locally on your device. Your data never leaves your computer.
              </Text>
            </Space>
          </div>
        }
        type="success"
        showIcon={false}
        style={{ border: 'none', background: 'transparent' }}
      />
    </div>
  );
});

PrivacyBanner.displayName = 'PrivacyBanner';

export default PrivacyBanner;
