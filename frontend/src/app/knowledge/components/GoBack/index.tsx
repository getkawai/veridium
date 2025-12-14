import { Icon } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { ArrowLeft } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';
import { useKnowledgeBaseStore } from '@/store/knowledgeBase';

const useStyles = createStyles(({ css, token }) => {
  return {
    container: css`
      cursor: pointer;
      color: ${token.colorTextDescription};

      &:hover {
        color: ${token.colorText};
      }
    `,
  };
});

interface GoBackProps {
  /**
   * The path to navigate to (relative to MemoryRouter)
   * e.g., "/" for /knowledge, "/bases" for /knowledge/bases
   */
  to?: string;
}

/**
 * GoBack component for react-router-dom
 * Uses useNavigate instead of Next.js Link
 */
const GoBack = memo<GoBackProps>(() => {
  const { t } = useTranslation('components');
  const { styles } = useStyles();
  const deactivateKnowledgeBase = useKnowledgeBaseStore((s) => s.deactivateKnowledgeBase);

  const handleClick = () => {
    deactivateKnowledgeBase();
  };

  return (
    <Flexbox align={'center'} className={styles.container} gap={4} horizontal onClick={handleClick}>
      <Icon icon={ArrowLeft} />
      <div>{t('GoBack.back')}</div>
    </Flexbox>
  );
});

GoBack.displayName = 'GoBack';

export default GoBack;
