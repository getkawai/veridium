import { Markdown } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { memo, useMemo } from 'react';

import { useIsMobile } from '@/hooks/useIsMobile';

import { useContainerStyles } from '../style';

const Preview = memo<{ content: string }>(({ content }) => {
  const { styles } = useContainerStyles();
  const isMobile = useIsMobile();

  // Custom components untuk desktop app link handling
  const markdownComponents = useMemo(
    () => ({
      a: ({ href, children, ...props }: any) => (
        <a
          {...props}
          href={href}
          onClick={(e) => {
            e.preventDefault();
            if (href) Browser.OpenURL(href);
          }}
          rel="noopener noreferrer"
          target="_blank"
        >
          {children}
        </a>
      ),
    }),
    []
  );

  return (
    <div className={styles.preview} style={{ padding: 12 }}>
      <Markdown 
        variant={isMobile ? 'chat' : undefined}
        components={markdownComponents} // Custom components untuk desktop link handling
      >
        {content}
      </Markdown>
    </div>
  );
});

export default Preview;
