import { ListLocalFileParams } from '@@/github.com/kawai-network/veridium/pkg/localfs';
import { ChatMessagePluginError } from '@/types';
import React, { memo } from 'react';

import { LocalFolder } from '@/features/LocalFile';
import { LocalFileListState } from '@@/github.com/kawai-network/veridium/pkg/yzma/tools/builtin';

import SearchResult from './Result';

interface ListFilesProps {
  args: ListLocalFileParams;
  messageId: string;
  pluginError: ChatMessagePluginError;
  pluginState?: LocalFileListState;
}

const ListFiles = memo<ListFilesProps>(({ messageId, pluginError, args, pluginState }) => {
  return (
    <>
      {args?.path && <LocalFolder path={args.path} />}
      <SearchResult
        listResults={pluginState?.listResults}
        messageId={messageId}
        pluginError={pluginError}
      />
    </>
  );
});

ListFiles.displayName = 'ListFiles';

export default ListFiles;
