import { createStyles } from 'antd-style';

/**
 * Shared styles untuk komponen ConfigPanel
 * Menyediakan animasi drag-and-drop yang konsisten
 */
export const useConfigPanelStyles = createStyles(({ css, token }) => ({
  // Style saat file di-drag over komponen
  dragOver: css`
    transform: scale(1.02);
    border-color: ${token.colorPrimary} !important;
    box-shadow: 0 0 0 2px ${token.colorPrimary}20;
    transition: transform 0.2s ease;
  `,

  // Transisi smooth untuk drag-and-drop
  dragTransition: css`
    transition:
      transform 0.2s ease,
      border-color 0.2s ease,
      box-shadow 0.2s ease;
  `,
}));
