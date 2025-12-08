'use client';

import { useAutoAnimate } from '@formkit/auto-animate/react';
import { ModelTag } from '@lobehub/icons';
import { ActionIconGroup, Block, Grid, Markdown, Tag, Text } from '@lobehub/ui';
import { App } from 'antd';
import { createStyles } from 'antd-style';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import { omit } from 'lodash-es';
import { CopyIcon, RotateCcwSquareIcon, Trash2 } from 'lucide-react';
import { RuntimeImageGenParams } from '@/model-bank';
import { memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

import InvalidAPIKey from '@/components/InvalidAPIKey';
import { useImageStore } from '@/store/image';
import { AsyncTaskErrorType } from '@/types/asyncTask';
import { GenerationBatch } from '@/types/generation';

import { GenerationItem } from './GenerationItem';
import { DEFAULT_MAX_ITEM_WIDTH } from './GenerationItem/utils';
import { ReferenceImages } from './ReferenceImages';

// =============================================================================
// STYLES DEFINITION
// =============================================================================

/**
 * Mendefinisikan styles menggunakan CSS-in-JS pattern dari antd-style
 * Menggunakan design tokens untuk konsistensi dengan tema aplikasi
 */
const useStyles = createStyles(({ cx, css, token }) => ({
  // Style untuk action buttons batch (copy, reuse, delete)
  // Secara default opacity 0 (tersembunyi), akan muncul saat hover container
  batchActions: cx(
    'batch-actions',
    css`
      opacity: 0;
      transition: opacity 0.1s ${token.motionEaseInOut};
    `,
  ),

  // Style khusus untuk tombol delete dengan warna merah saat hover
  batchDeleteButton: css`
    &:hover {
      border-color: ${token.colorError} !important;
      color: ${token.colorError} !important;
      background: ${token.colorErrorBg} !important;
    }
  `,

  // Style container utama yang menampilkan action buttons saat di-hover
  container: css`
    &:hover {
      .batch-actions {
        opacity: 1;
      }
    }
  `,

  // Style untuk area prompt dengan custom styling untuk <pre> tags
  prompt: css`
    pre {
      overflow: hidden !important;
      padding-block: 4px;
      font-size: 13px;
    }
  `,
}));

// Extend dayjs dengan plugin relativeTime untuk menampilkan waktu relatif
// (contoh: "2 hours ago", "yesterday")
dayjs.extend(relativeTime);

// =============================================================================
// TYPE DEFINITIONS
// =============================================================================

/**
 * Props interface untuk komponen GenerationBatchItem
 * @property batch - Data batch generasi gambar yang akan ditampilkan
 */
interface GenerationBatchItemProps {
  batch: GenerationBatch;
}

// =============================================================================
// MAIN COMPONENT
// =============================================================================

/**
 * GenerationBatchItem - Komponen untuk menampilkan satu batch generasi gambar
 *
 * Fitur utama:
 * - Menampilkan prompt yang digunakan untuk generasi
 * - Menampilkan metadata (model, ukuran, jumlah gambar)
 * - Menampilkan grid gambar hasil generasi
 * - Menampilkan gambar referensi (jika ada)
 * - Action buttons: copy prompt, reuse settings, delete batch
 * - Handle error API key tidak valid
 *
 * Komponen ini di-memo untuk optimasi performa (mencegah re-render yang tidak perlu)
 */
export const GenerationBatchItem = memo<GenerationBatchItemProps>(({ batch }) => {
  // ==========================================================================
  // HOOKS INITIALIZATION
  // ==========================================================================

  const { styles } = useStyles();
  const { t } = useTranslation(['image', 'modelProvider', 'error']);
  const { message } = App.useApp();

  // Hook untuk auto-animate pada grid gambar
  // Memberikan animasi otomatis saat item ditambah/dihapus dari grid
  const [imageGridRef] = useAutoAnimate();

  // ==========================================================================
  // STORE SELECTORS
  // ==========================================================================

  const activeTopicId = useImageStore((s) => s.activeGenerationTopicId);
  const removeGenerationBatch = useImageStore((s) => s.removeGenerationBatch);
  const recreateImage = useImageStore((s) => s.recreateImage);
  const reuseSettings = useImageStore((s) => s.reuseSettings);

  // ==========================================================================
  // MEMOIZED VALUES
  // ==========================================================================

  /**
   * Format waktu pembuatan batch ke format yang lebih readable
   * Di-memo untuk menghindari re-calculation yang tidak perlu
   */
  const time = useMemo(() => {
    return dayjs(batch.createdAt).format('YYYY-MM-DD HH:mm:ss');
  }, [batch.createdAt]);

  // ==========================================================================
  // EVENT HANDLERS
  // ==========================================================================

  /**
   * Handler untuk menyalin prompt ke clipboard
   * Menampilkan toast success/error berdasarkan hasil operasi
   */
  const handleCopyPrompt = async () => {
    try {
      await navigator.clipboard.writeText(batch.prompt);
      message.success(t('generation.actions.promptCopied'));
    } catch (error) {
      console.error('Failed to copy prompt:', error);
      message.error(t('generation.actions.promptCopyFailed'));
    }
  };

  /**
   * Handler untuk menggunakan ulang settings dari batch ini
   * Mengambil model, provider, dan config (tanpa seed) untuk digunakan di generasi baru
   */
  const handleReuseSettings = () => {
    reuseSettings(
      batch.model,
      batch.provider,
      // Menghilangkan seed agar generasi baru menggunakan seed random
      omit(batch.config as RuntimeImageGenParams, ['seed']),
    );
  };

  /**
   * Handler untuk menghapus batch
   * Memerlukan activeTopicId untuk mengetahui batch milik topik mana
   */
  const handleDeleteBatch = async () => {
    // Guard clause: tidak melakukan apa-apa jika tidak ada topik aktif
    if (!activeTopicId) return;

    try {
      await removeGenerationBatch(batch.id, activeTopicId);
    } catch (error) {
      console.error('Failed to delete batch:', error);
    }
  };

  // ==========================================================================
  // EARLY RETURNS (Guard Clauses)
  // ==========================================================================

  // Jika batch tidak memiliki generasi, tidak render apa-apa
  if (batch.generations.length === 0) {
    return null;
  }

  // ==========================================================================
  // ERROR HANDLING - Invalid API Key
  // ==========================================================================

  /**
   * Cek apakah ada error API key tidak valid di salah satu generasi
   * Menggunakan Array.some() untuk early exit jika ditemukan
   */
  const isInvalidApiKey = batch.generations.some(
    (generation) => generation.task.error?.name === AsyncTaskErrorType.InvalidProviderAPIKey,
  );

  // Jika API key tidak valid, tampilkan komponen error khusus
  if (isInvalidApiKey) {
    return (
      <InvalidAPIKey
        bedrockDescription={t('bedrock.unlock.imageGenerationDescription', { ns: 'modelProvider' })}
        description={t('unlock.apiKey.imageGenerationDescription', {
          name: batch.provider,
          ns: 'error',
        })}
        id={batch.id}
        onClose={() => {
          removeGenerationBatch(batch.id, activeTopicId!);
        }}
        onRecreate={() => {
          recreateImage(batch.id);
        }}
        provider={batch.provider}
      />
    );
  }

  // ==========================================================================
  // LAYOUT CALCULATIONS
  // ==========================================================================

  /**
   * Menghitung total jumlah gambar referensi
   * - imageUrl: single image reference (0 atau 1)
   * - imageUrls: multiple image references (array)
   */
  const referenceImageCount =
    (batch.config?.imageUrl ? 1 : 0) + (batch.config?.imageUrls?.length || 0);

  // Menentukan apakah menggunakan layout single image
  // Layout berbeda untuk 1 gambar referensi vs multiple/no gambar
  const isSingleImageLayout = referenceImageCount === 1;

  // ==========================================================================
  // REUSABLE UI ELEMENTS
  // ==========================================================================

  /**
   * Komponen untuk menampilkan prompt dan metadata
   * Di-extract menjadi variable untuk digunakan di kedua layout (single/multiple)
   */
  const promptAndMetadata = (
    <>
      {/* Menampilkan prompt dengan format Markdown */}
      <Markdown variant={'chat'}>{batch.prompt}</Markdown>

      {/* Container untuk metadata tags */}
      <Flexbox gap={4} horizontal justify="space-between" style={{ marginBottom: 10 }}>
        <Flexbox gap={4} horizontal>
          {/* Tag untuk menampilkan model AI yang digunakan */}
          <ModelTag model={batch.model} />

          {/* Tag untuk ukuran gambar (jika tersedia) */}
          {batch.width && batch.height && (
            <Tag>
              {batch.width} × {batch.height}
            </Tag>
          )}

          {/* Tag untuk jumlah gambar yang di-generate */}
          <Tag>{t('generation.metadata.count', { count: batch.generations.length })}</Tag>
        </Flexbox>
      </Flexbox>
    </>
  );

  // ==========================================================================
  // RENDER
  // ==========================================================================

  return (
    <Block className={styles.container} gap={8} variant="borderless">
      {/* 
        Conditional Layout berdasarkan jumlah gambar referensi:
        - Single image: layout horizontal dengan vertical centering
        - Multiple/no images: layout vertical
      */}
      {isSingleImageLayout ? (
        // ===== SINGLE IMAGE LAYOUT =====
        // Gambar referensi di kiri, prompt dan metadata di kanan
        <Flexbox align="center" gap={16} horizontal>
          <ReferenceImages
            imageUrl={batch.config?.imageUrl}
            imageUrls={batch.config?.imageUrls}
            layout="single"
          />
          <Flexbox flex={1} gap={8}>
            {promptAndMetadata}
          </Flexbox>
        </Flexbox>
      ) : (
        // ===== MULTIPLE/NO IMAGES LAYOUT =====
        // Gambar referensi di atas, prompt dan metadata di bawah
        <>
          <ReferenceImages
            imageUrl={batch.config?.imageUrl}
            imageUrls={batch.config?.imageUrls}
            layout="multiple"
          />
          {promptAndMetadata}
        </>
      )}

      {/* 
        Grid untuk menampilkan gambar-gambar hasil generasi
        - maxItemWidth: batasi lebar maksimum setiap item
        - ref: untuk auto-animate
        - rows: jumlah baris sesuai jumlah generasi
      */}
      <Grid
        maxItemWidth={DEFAULT_MAX_ITEM_WIDTH}
        ref={imageGridRef}
        rows={batch.generations.length}
      >
        {/* Iterasi setiap generasi dan render GenerationItem */}
        {batch.generations.map((generation) => (
          <GenerationItem
            generation={generation}
            generationBatch={batch}
            key={generation.id}
            prompt={batch.prompt}
          />
        ))}
      </Grid>

      {/* 
        Action buttons container
        - Tersembunyi secara default (opacity: 0)
        - Muncul saat container di-hover
      */}
      <Flexbox
        align={'center'}
        className={styles.batchActions}
        horizontal
        justify={'space-between'}
      >
        {/* Timestamp pembuatan batch */}
        <Text as={'time'} fontSize={12} type={'secondary'}>
          {time}
        </Text>

        {/* 
          Group of action buttons:
          1. Reuse Settings - menggunakan ulang settings untuk generasi baru
          2. Copy Prompt - menyalin prompt ke clipboard
          3. Delete Batch - menghapus batch (dengan styling danger/merah)
        */}
        <ActionIconGroup
          items={[
            {
              icon: RotateCcwSquareIcon,
              key: 'reuseSettings',
              label: t('generation.actions.reuseSettings'),
              onClick: handleReuseSettings,
            },
            {
              icon: CopyIcon,
              key: 'copyPrompt',
              label: t('generation.actions.copyPrompt'),
              onClick: handleCopyPrompt,
            },
            {
              danger: true,
              icon: Trash2,
              key: 'deleteBatch',
              label: t('generation.actions.deleteBatch'),
              onClick: handleDeleteBatch,
            },
          ]}
        />
      </Flexbox>
    </Block>
  );
});

// Display name untuk debugging di React DevTools
GenerationBatchItem.displayName = 'GenerationBatchItem';
