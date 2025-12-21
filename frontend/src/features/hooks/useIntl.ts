import { useTranslation } from 'react-i18next';

const useIntl = (componentName: string) => {
  const { t } = useTranslation('common'); // Fallback to common or specific namespace if available

  // Mocking the structure expected by PayPanel
  // In a real scenario, these should be added to the locale files
  const messages = {
    tips: t('payPanel.tips', 'Please scan the QR code to pay'),
  };

  return { messages };
};

export default useIntl;
