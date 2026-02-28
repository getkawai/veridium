import { MoveLocalFilesParams } from '@@/github.com/getkawai/tools/localfs';
import { ChatMessagePluginError } from '@/types';
import { Icon } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { ArrowRightIcon, CheckCircle2, XCircle } from 'lucide-react';
import path from 'path-browserify-esm';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

import { LocalFile } from '@/features/LocalFile';
import { LocalMoveFilesState } from '../../type';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    color: ${token.colorTextQuaternary};
  `,
  item: css`
    padding: 4px 0;
  `,
  success: css`
    color: ${token.colorSuccess};
  `,
  error: css`
    color: ${token.colorError};
  `,
  summary: css`
    font-size: 12px;
    color: ${token.colorTextSecondary};
    margin-top: 8px;
  `,
}));

interface MoveFilesProps {
  args: MoveLocalFilesParams;
  messageId: string;
  pluginError: ChatMessagePluginError;
  pluginState: LocalMoveFilesState;
}

const MoveFiles = memo<MoveFilesProps>(({ args, pluginState }) => {
  const { styles } = useStyles();

  if (!pluginState?.results?.length && !args?.items?.length) return null;

  const results = pluginState?.results || [];
  const items = args?.items || [];

  return (
    <Flexbox className={styles.container} gap={4} paddingInline={12}>
      {results.length > 0 ? (
        <>
          {results.map((result, index) => {
            const oldFileName = path.basename(result.sourcePath || '');
            const newFileName = path.basename(result.newPath || '');
            
            return (
              <Flexbox 
                key={index} 
                align={'center'} 
                className={styles.item} 
                gap={8} 
                horizontal
              >
                <Icon 
                  icon={result.success ? CheckCircle2 : XCircle} 
                  className={result.success ? styles.success : styles.error}
                  size={14}
                />
                <Flexbox>{oldFileName}</Flexbox>
                <Icon icon={ArrowRightIcon} size={12} />
                <LocalFile name={newFileName} path={result.newPath} />
              </Flexbox>
            );
          })}
          {pluginState.totalCount > 0 && (
            <Flexbox className={styles.summary}>
              {pluginState.successCount}/{pluginState.totalCount} files moved successfully
            </Flexbox>
          )}
        </>
      ) : (
        items.map((item, index) => {
          const oldFileName = path.basename(item.oldPath || '');
          const newFileName = path.basename(item.newPath || '');
          
          return (
            <Flexbox 
              key={index} 
              align={'center'} 
              className={styles.item} 
              gap={8} 
              horizontal
            >
              <Flexbox>{oldFileName}</Flexbox>
              <Icon icon={ArrowRightIcon} size={12} />
              <LocalFile name={newFileName} path={item.newPath} />
            </Flexbox>
          );
        })
      )}
    </Flexbox>
  );
});

export default MoveFiles;
