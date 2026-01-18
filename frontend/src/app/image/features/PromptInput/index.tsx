'use client';

import { Button, TextArea } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { Sparkles } from 'lucide-react';
import type { KeyboardEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import { useGeminiChineseWarning } from '@/hooks/useGeminiChineseWarning';
import { useImageStore } from '@/store/image';
import { createImageSelectors } from '@/store/image/selectors';
import { useGenerationConfigParam } from '@/store/image/slices/generationConfig/hooks';
import { imageGenerationConfigSelectors } from '@/store/image/slices/generationConfig/selectors';

import PromptTitle from './Title';

// =============================================================================
// TYPE DEFINITIONS
// =============================================================================

/**
 * Props interface untuk komponen PromptInput
 * @property disableAnimation - Opsional, untuk menonaktifkan animasi
 * @property showTitle - Opsional, menampilkan judul di atas input (default: false)
 */
interface PromptInputProps {
  disableAnimation?: boolean;
  showTitle?: boolean;
}

// =============================================================================
// STYLES DEFINITION
// =============================================================================

/**
 * Mendefinisikan styles untuk komponen PromptInput
 * Menggunakan design tokens untuk konsistensi dengan tema aplikasi
 */
const useStyles = createStyles(({ css, token, isDarkMode }) => ({
  // Style untuk container input prompt
  // Memiliki border, border-radius, background, dan shadow
  // Shadow berbeda untuk dark mode dan light mode
  container: css`
    border: 1px solid ${token.colorBorderSecondary};
    border-radius: ${token.borderRadiusLG * 1.5}px;
    background-color: ${token.colorBgContainer};
    box-shadow:
      ${token.boxShadowTertiary},
      ${isDarkMode
        ? `0 0 48px 32px ${token.colorBgContainerSecondary}`
        : `0 0 0  ${token.colorBgContainerSecondary}`},
      0 32px 0 ${token.colorBgContainerSecondary};
  `,

  // Style wrapper untuk layout flexbox dengan gap dan alignment
  wrapper: css`
    display: flex;
    flex-direction: column;
    gap: 16px;
    align-items: center;

    width: 100%;
  `,
}));

// =============================================================================
// MAIN COMPONENT
// =============================================================================

/**
 * PromptInput - Komponen input untuk memasukkan prompt generasi gambar
 *
 * Fitur utama:
 * - TextArea untuk input prompt multi-line
 * - Tombol generate dengan icon sparkles
 * - Validasi login sebelum generate
 * - Warning untuk penggunaan bahasa Cina dengan model Gemini
 * - Keyboard shortcut: Enter untuk generate (Shift+Enter untuk new line)
 * - Loading state saat proses generasi berlangsung
 */
const PromptInput = ({ showTitle = false }: PromptInputProps) => {
  // ==========================================================================
  // HOOKS INITIALIZATION
  // ==========================================================================

  const { styles } = useStyles();
  const { t } = useTranslation('image');

  // Hook untuk mengambil dan mengubah nilai prompt dari store
  const { value, setValue } = useGenerationConfigParam('prompt');

  // ==========================================================================
  // STORE SELECTORS
  // ==========================================================================

  // Status apakah sedang dalam proses membuat gambar
  const isCreating = useImageStore(createImageSelectors.isCreating);

  // Action untuk membuat/generate gambar
  const createImage = useImageStore((s) => s.createImage);

  // Model AI yang sedang dipilih untuk generasi
  const currentModel = useImageStore(imageGenerationConfigSelectors.model);

  // Hook untuk menampilkan warning jika menggunakan bahasa Cina dengan Gemini
  const checkGeminiChineseWarning = useGeminiChineseWarning();

  // ==========================================================================
  // EVENT HANDLERS
  // ==========================================================================

  /**
   * Handler untuk memulai proses generasi gambar
   * Melakukan validasi:
   * 1. User harus login
   * 2. Cek warning bahasa Cina untuk model Gemini
   * Jika validasi lolos, panggil createImage()
   */
  const handleGenerate = async () => {
    // Cek dan tampilkan warning jika menggunakan teks Cina dengan model Gemini
    const shouldContinue = await checkGeminiChineseWarning({
      model: currentModel,
      prompt: value,
      scenario: 'image',
    });

    // Batalkan jika user memilih untuk tidak melanjutkan
    if (!shouldContinue) return;

    // Jalankan proses generasi gambar
    await createImage();
  };

  /**
   * Handler untuk keyboard event pada TextArea
   * Enter tanpa Shift = submit/generate
   * Shift+Enter = new line (default behavior)
   */
  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      // Hanya generate jika tidak sedang loading dan prompt tidak kosong
      if (!isCreating && value.trim()) {
        handleGenerate();
      }
    }
  };

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <Flexbox
      gap={32}
      style={{
        marginTop: 48,
      }}
      width={'100%'}
    >
      {/* Tampilkan judul jika showTitle = true */}
      {showTitle && <PromptTitle />}

      {/* Container utama untuk input dan tombol generate */}
      <Flexbox
        align="flex-end"
        className={styles.container}
        gap={12}
        height={'100%'}
        horizontal
        padding={'12px 12px 12px 16px'}
        width={'100%'}
      >
        {/* 
          TextArea untuk input prompt
          - autoSize: otomatis resize antara 3-6 baris
          - variant borderless: tanpa border default
        */}
        <TextArea
          autoSize={{ maxRows: 6, minRows: 3 }}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={t('config.prompt.placeholder')}
          style={{
            borderRadius: 0,
            padding: 0,
          }}
          value={value}
          variant={'borderless'}
        />

        {/* 
          Tombol Generate
          - disabled: jika prompt kosong
          - loading: menampilkan spinner saat proses generasi
          - icon Sparkles: menunjukkan fitur AI generation
        */}
        <Button
          disabled={!value}
          icon={Sparkles}
          loading={isCreating}
          onClick={handleGenerate}
          size={'large'}
          style={{
            fontWeight: 500,
            height: 64,
            minWidth: 64,
            width: 64,
          }}
          title={isCreating ? t('generation.status.generating') : t('generation.actions.generate')}
          type={'primary'}
        />
      </Flexbox>
    </Flexbox>
  );
};

export default PromptInput;
