import { Icon, ItemType } from '@lobehub/ui';
import isEqual from 'fast-deep-equal';
import { ArrowRight, LibraryBig } from 'lucide-react';
import { memo, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import FileIcon from '@/components/FileIcon';
import RepoIcon from '@/components/RepoIcon';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/selectors';
import { createServiceLogger } from '@/utils/logger';

const logger = createServiceLogger('ChatInput', 'useControls', 'features/ChatInput/ActionBar/Knowledge/useControls.tsx');

import CheckboxItem from '../components/CheckbokWithLoading';

export const useControls = ({
  setModalOpen,
}: {
  setModalOpen: (open: boolean) => void;
}) => {
  /* eslint-disable sort-keys-fix/sort-keys-fix */
  const { internal_refreshAgentConfig, activeId } = useAgentStore();
  const agentConfig = useAgentStore(agentSelectors.currentAgentConfig, isEqual);
  const [updating, setUpdating] = useState(false);
  const { t } = useTranslation('chat');

  const files = useAgentStore(agentSelectors.currentAgentFiles, isEqual);
  const knowledgeBases = useAgentStore(agentSelectors.currentAgentKnowledgeBases, isEqual);
  logger.info('[useControls] Files from store:', files);

  // Self-healing: if we have an agent but no files loaded (and we expect files might exist, or just to be safe),
  // trigger a refresh once on mount/agent change.
  // Ideally we check if "loaded" status is true, but we can just trigger a refresh to be safe if list is empty.
  useEffect(() => {
    if (activeId && agentConfig?.id && files.length === 0 && knowledgeBases.length === 0) {
      logger.info('[useControls] Empty files/KBs detected, triggering refresh for', activeId);
      internal_refreshAgentConfig(activeId);
    }
  }, [activeId, agentConfig?.id, files.length, knowledgeBases.length]);

  const [toggleFile, toggleKnowledgeBase] = useAgentStore((s) => [
    s.toggleFile,
    s.toggleKnowledgeBase,
  ]);

  const items: ItemType[] = [
    // {
    //   children: [
    //     {
    //       icon: <RepoIcon />,
    //       key: 'allFiles',
    //       label: <KnowledgeBaseItem id={'all'} label={t('knowledgeBase.allFiles')} />,
    //     },
    //     {
    //       icon: <RepoIcon />,
    //       key: 'allRepos',
    //       label: <KnowledgeBaseItem id={'all'} label={t('knowledgeBase.allKnowledgeBases')} />,
    //     },
    //   ],
    //   key: 'all',
    //   label: (
    //     <Flexbox horizontal justify={'space-between'}>
    //       {t('knowledgeBase.all')}
    //       {/*<Link href={'/files'}>{t('knowledgeBase.more')}</Link>*/}
    //     </Flexbox>
    //   ),
    //   type: 'group',
    // },
    {
      children: [
        // first the files
        ...files.map((item) => ({
          icon: <FileIcon fileName={item.name} fileType={item.type} size={20} />,
          key: item.id,
          label: (
            <CheckboxItem
              checked={item.enabled}
              id={item.id}
              label={item.name}
              onUpdate={async (id, enabled) => {
                setUpdating(true);
                await toggleFile(id, enabled);
                setUpdating(false);
              }}
            />
          ),
        })),

        // then the knowledge bases
        ...knowledgeBases.map((item) => ({
          icon: <RepoIcon />,
          key: item.id,
          label: (
            <CheckboxItem
              checked={item.enabled}
              id={item.id}
              label={item.name}
              onUpdate={async (id, enabled) => {
                setUpdating(true);
                await toggleKnowledgeBase(id, enabled);
                setUpdating(false);
              }}
            />
          ),
        })),
      ],
      key: 'relativeFilesOrKnowledgeBases',
      label: t('knowledgeBase.relativeFilesOrKnowledgeBases'),
      type: 'group',
    },
    {
      type: 'divider',
    },
    {
      extra: <Icon icon={ArrowRight} />,
      icon: LibraryBig,
      key: 'knowledge-base-store',
      label: t('knowledgeBase.viewMore'),
      onClick: () => {
        setModalOpen(true);
      },
    },
  ];

  return items;
};
