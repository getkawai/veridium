import { App } from 'antd';
import { Copy } from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';

interface CopyButtonProps {
  text: string;
}

export const CopyButton = ({ text }: CopyButtonProps) => {
  const { message } = App.useApp();

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text);
      message.success("Copied!");
    } catch (err) {
      console.error('Failed to copy:', err);
      message.error("Failed to copy. Please try again.");
    }
  };

  return (
    <ActionIcon
      icon={Copy}
      size="small"
      onClick={handleCopy}
      title="Copy"
    />
  );
};
