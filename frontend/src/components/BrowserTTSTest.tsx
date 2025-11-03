/**
 * Test component untuk Browser TTS
 * Dapat digunakan untuk testing cepat fitur Browser TTS
 * 
 * Usage:
 * import BrowserTTSTest from '@/components/BrowserTTSTest';
 * <BrowserTTSTest />
 */

import { Button, Select } from '@lobehub/ui';
import { useEffect, useState } from 'react';
import { Flexbox } from 'react-layout-kit';

import { useBrowserVoices } from '@/hooks/useBrowserVoices';
import { getBrowserVoices, isBrowserTTSSupported } from '@/utils/browserTTS';

export const BrowserTTSTest = () => {
  const [selectedVoice, setSelectedVoice] = useState<string>('');
  const [testText, setTestText] = useState('Hello, this is a test of browser text to speech.');
  const [isPlaying, setIsPlaying] = useState(false);
  const { voices, loading } = useBrowserVoices();

  useEffect(() => {
    if (voices.length > 0 && !selectedVoice) {
      setSelectedVoice(voices[0].value);
    }
  }, [voices, selectedVoice]);

  const handleSpeak = () => {
    if (!window.speechSynthesis || !selectedVoice) return;

    window.speechSynthesis.cancel();

    const utterance = new SpeechSynthesisUtterance(testText);
    const allVoices = window.speechSynthesis.getVoices();
    const voice = allVoices.find((v) => v.name === selectedVoice);

    if (voice) {
      utterance.voice = voice;
    }

    utterance.onstart = () => setIsPlaying(true);
    utterance.onend = () => setIsPlaying(false);
    utterance.onerror = () => setIsPlaying(false);

    window.speechSynthesis.speak(utterance);
  };

  const handleStop = () => {
    window.speechSynthesis.cancel();
    setIsPlaying(false);
  };

  const handleGetVoices = async () => {
    const browserVoices = await getBrowserVoices();
    console.log('Available browser voices:', browserVoices);
  };

  if (!isBrowserTTSSupported()) {
    return (
      <Flexbox padding={16}>
        <div>Browser TTS is not supported in this browser</div>
      </Flexbox>
    );
  }

  return (
    <Flexbox gap={16} padding={16}>
      <h2>Browser TTS Test</h2>

      <div>
        <label>Test Text:</label>
        <textarea
          onChange={(e) => setTestText(e.target.value)}
          rows={3}
          style={{ width: '100%' }}
          value={testText}
        />
      </div>

      <div>
        <label>Select Voice:</label>
        <Select
          loading={loading}
          onChange={(value) => setSelectedVoice(value as string)}
          options={voices}
          placeholder="Select a voice"
          style={{ width: '100%' }}
          value={selectedVoice}
        />
      </div>

      <Flexbox gap={8} horizontal>
        <Button disabled={isPlaying || !selectedVoice} onClick={handleSpeak} type="primary">
          {isPlaying ? 'Playing...' : 'Speak'}
        </Button>
        <Button disabled={!isPlaying} onClick={handleStop}>
          Stop
        </Button>
        <Button onClick={handleGetVoices}>Log Voices to Console</Button>
      </Flexbox>

      <div style={{ fontSize: '12px', color: '#666' }}>
        Available voices: {voices.length}
      </div>
    </Flexbox>
  );
};

export default BrowserTTSTest;

