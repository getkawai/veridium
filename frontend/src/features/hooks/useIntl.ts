import { useTranslation } from 'react-i18next';

const useIntl = (ns: string = 'payPanel') => {
  const { t } = useTranslation(ns);

  // Bridging the structure expected by PayPanel
  // This allows the component to access translations via messages.key
  const proxyHandler = {
    get: (target: any, prop: string) => {
      return t(prop);
    },
  };

  const messages = new Proxy({}, proxyHandler) as Record<string, string>;

  return { messages, t };
};

export default useIntl;
