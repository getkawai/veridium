import { LocalReadFileParams } from '@@/github.com/getkawai/tools/localfs';
import { ChatMessagePluginError } from '@/types';
import { memo } from 'react';

import { useChatStore } from '@/store/chat';
import { chatToolSelectors } from '@/store/chat/slices/builtinTool/selectors';
import { LocalReadFileState } from '@@/github.com/getkawai/tools/builtin';

import ReadFileSkeleton from './ReadFileSkeleton';
import ReadFileView from './ReadFileView';

interface ReadFileQueryProps {
  args: LocalReadFileParams;
  messageId: string;
  pluginError: ChatMessagePluginError;
  pluginState: LocalReadFileState;
}

const ReadFileQuery = memo<ReadFileQueryProps>(({ args, pluginState, messageId }) => {
  const loading = useChatStore(chatToolSelectors.isSearchingLocalFiles(messageId));

  if (loading) {
    return <ReadFileSkeleton />;
  }

  if (!args?.path || !pluginState) return null;

  return (
    <ReadFileView
      {...pluginState.fileContent}
      loc={pluginState.fileContent.loc as [number, number]}
      path={args.path}
    />
  );
});

export default ReadFileQuery;
