'use client';

import { useCallback } from 'react';

import { useSetFileModalId } from '../../shared/useFileQueryParam';
import FileDetail from './FileDetail';
import FilePreview from './FilePreview';
import FullscreenModal from './FullscreenModal';

interface ModalPageClientProps {
  id: string;
}

const ModalPageClient = ({ id }: ModalPageClientProps) => {
  const setFileModalId = useSetFileModalId();

  const handleClose = useCallback(() => {
    setFileModalId(undefined);
  }, [setFileModalId]);

  return (
    <FullscreenModal detail={<FileDetail id={id} />} onClose={handleClose}>
      <FilePreview id={id} />
    </FullscreenModal>
  );
};

export default ModalPageClient;
