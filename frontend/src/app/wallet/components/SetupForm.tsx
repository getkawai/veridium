import { memo, useEffect, useState } from 'react';
import { Form, Input, Button } from 'antd';
import { Flexbox } from 'react-layout-kit';
import { App } from 'antd';
import { useUserStore } from '@/store/user';
import { CopyButton } from './CopyButton';

interface SetupFormProps {
  type: 'create' | 'import';
  onSuccess: () => void;
}

export const SetupForm = memo<SetupFormProps>(({ type, onSuccess }) => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [mnemonic, setMnemonic] = useState('');
  const [description, setDescription] = useState('');
  const [step, setStep] = useState<'form' | 'mnemonic'>(type === 'create' ? 'mnemonic' : 'form');
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();

  const { generateMnemonic, createWallet } = useUserStore();

  useEffect(() => {
    if (type === 'create' && step === 'mnemonic') {
      generateMnemonic().then(setMnemonic);
    }
  }, [type, step, generateMnemonic]);

  const handleFinish = async () => {
    if (password !== confirmPassword) {
      return message.error("Passwords do not match");
    }
    if (password.length < 8) {
      return message.error("Password too short");
    }
    setLoading(true);
    try {
      await createWallet(password, mnemonic, description);
      message.success("Wallet created successfully");
      onSuccess();
    } catch (e: any) {
      message.error(e.message || "Failed to create wallet");
    } finally {
      setLoading(false);
    }
  };

  if (type === 'create' && step === 'mnemonic') {
    return (
      <Flexbox gap={12}>
        <p>Save these 12 words securely:</p>
        <div style={{ background: 'rgba(0,0,0,0.05)', padding: 16, borderRadius: 8, textAlign: 'center' }}>
          <code style={{ fontSize: 16, fontWeight: 700 }}>{mnemonic}</code>
          <div style={{ marginTop: 8 }}>
            <CopyButton text={mnemonic} />
          </div>
        </div>
        <Button type="primary" block onClick={() => setStep('form')}>I have written it down</Button>
      </Flexbox>
    );
  }

  return (
    <Form layout="vertical" onFinish={handleFinish}>
      {type === 'import' && (
        <Form.Item label="Mnemonic Phrase" required>
          <Input.TextArea
            rows={3}
            value={mnemonic}
            onChange={e => setMnemonic(e.target.value)}
            placeholder="word1 word2 ..."
          />
        </Form.Item>
      )}
      <Form.Item label="Account Description">
        <Input
          placeholder="Main account, Savings, etc."
          value={description}
          onChange={e => setDescription(e.target.value)}
        />
      </Form.Item>
      <Form.Item label="Lock Password" required>
        <Input.Password
          value={password}
          onChange={e => setPassword(e.target.value)}
          placeholder="At least 8 characters"
        />
      </Form.Item>
      <Form.Item label="Confirm Password" required>
        <Input.Password
          value={confirmPassword}
          onChange={e => setConfirmPassword(e.target.value)}
          placeholder="Repeat password"
        />
      </Form.Item>
      <Button type="primary" htmlType="submit" block size="large" loading={loading}>
        {type === 'create' ? 'Create Account' : 'Import Account'}
      </Button>
    </Form>
  );
});

SetupForm.displayName = 'SetupForm';
