'use client';

import { Text } from '@lobehub/ui';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import FileList from './FileList';
import Header from './Header';
import UploadDock from './UploadDock';
import ChunkDrawer from './ChunkDrawer';

interface FileManagerProps {
  category?: string;
  knowledgeBaseId?: string;
  onOpenFile: (id: string) => void;
  title: string;
}
const FileManager = memo<FileManagerProps>(({ title, knowledgeBaseId, category, onOpenFile }) => {
  return (
    <>
      <Header knowledgeBaseId={knowledgeBaseId} />
      <Flexbox gap={12} height={'100%'}>
        <Text strong style={{ fontSize: 16, marginBlock: 16, marginInline: 24 }}>
          {title}
        </Text>
        <FileList category={category} knowledgeBaseId={knowledgeBaseId} onOpenFile={onOpenFile} />
      </Flexbox>
      <UploadDock />
      <ChunkDrawer />
    </>
  );
});

export default FileManager;
