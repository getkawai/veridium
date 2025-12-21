import { createStyles } from 'antd-style';

export const useStyles = createStyles(({ css, token, prefixCls }) => {
  const componentCls = `${prefixCls}-web3-pay-panel`;

  return {
    container: css`
      &.${componentCls} {
        .${componentCls}-content {
          .${componentCls}-title {
            font-size: ${token.fontSizeHeading4}px;
            line-height: ${token.lineHeightHeading4};
            color: ${token.colorTextBase};
            padding-bottom: ${token.paddingLG}px;
          }
          .${componentCls}-desc {
            font-size: ${token.fontSize}px;
            color: ${token.colorTextSecondary};
            line-height: ${token.lineHeightSM};
          }
          .${componentCls}-chainItem {
            width: 100%;
            padding-block: ${token.paddingXS}px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            cursor: pointer;
            .${componentCls}-chainInfo {
              font-size: ${token.fontSizeHeading5}px;
              display: flex;
              align-items: center;
              .${componentCls}-icon {
                font-size: ${token.fontSizeHeading1}px;
                padding-inline-end: ${token.paddingSM}px;
              }
              .${componentCls}-type {
                font-size: ${token.fontSize}px;
                color: ${token.colorTextDescription};
                line-height: ${token.lineHeightSM};
              }
            }
            .${componentCls}-gasInfo {
              font-size: ${token.fontSize}px;
              color: ${token.colorTextDescription};
              line-height: ${token.lineHeightSM};
            }
          }
          .${componentCls}-code-title {
            font-size: ${token.fontSize}px;
            line-height: ${token.lineHeightHeading4};
            color: ${token.colorTextBase};
            padding-bottom: ${token.padding}px;
          }
          .${componentCls}-code-amount {
            font-size: ${token.fontSizeHeading1}px;
            line-height: ${token.lineHeightHeading1};
            color: ${token.colorTextBase};
            padding-block: ${token.padding}px;
          }

          .${componentCls}-code-content {
            display: flex;
            flex-direction: column;
            align-items: center;
            .${componentCls}-code-desc {
              font-size: ${token.fontSize}px;
              color: ${token.colorTextSecondary};
              line-height: ${token.lineHeightSM};
              width: 100%;
            }
            .${componentCls}-code-tips {
              font-size: ${token.fontSize}px;
              line-height: ${token.lineHeightHeading1};
              padding-block: ${token.padding}px;
              color: ${token.colorTextDescription};
              display: flex;
              gap: ${token.paddingSM}px;
            }
          }
        }
      }
    `,
  };
});
