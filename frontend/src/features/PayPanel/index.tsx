import React, { useContext, useState } from 'react';
import { ConfigProvider } from 'antd';
import classNames from 'classnames';

import { ChainList } from './ChainList';
import { PayPanelContext, type PayPanelProps } from './PayPanelContext';
import { ShowCode } from './ShowCode';
import { useStyles } from './style';

export const PayPanel: React.FC<React.PropsWithChildren<PayPanelProps>> = (props) => {
  const { getPrefixCls } = useContext(ConfigProvider.ConfigContext);
  const prefixCls = getPrefixCls('web3-pay-panel');
  const { styles } = useStyles();
  const [selectedChainId, setSelectedChainId] = useState<number | string>();

  const { amount, target, supportedChains, token, wallets, onFinish } = props;

  return (
    <PayPanelContext.Provider
      value={{
        amount,
        target,
        supportedChains,
        token,
        wallets,
        onFinish,
      }}
    >
      <div className={classNames(prefixCls, styles.container)}>
        <div className={`${prefixCls}-content`}>
          {!selectedChainId && (
            <ChainList
              onChainSelected={(e) => {
                setSelectedChainId(e);
              }}
            />
          )}
          {selectedChainId && (
            <ShowCode
              selectedChainId={selectedChainId}
              onReturn={() => {
                setSelectedChainId(undefined);
              }}
            />
          )}
        </div>
      </div>
    </PayPanelContext.Provider>
  );
};
