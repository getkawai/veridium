import { Markdown, Mermaid } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { memo, lazy, useMemo } from 'react';

import HTMLRenderer from './HTML';
import SVGRender from './SVG';

const ReactRenderer = lazy(() => import('./React'));

const Renderer = memo<{ content: string; type?: string }>(({ content, type }) => {
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

  switch (type) {
    case 'application/lobe.artifacts.react': {
      return <ReactRenderer code={content} />;
    }

    case 'image/svg+xml': {
      return <SVGRender content={content} />;
    }

    case 'application/lobe.artifacts.mermaid': {
      return <Mermaid variant={'borderless'}>{content}</Mermaid>;
    }

    case 'text/markdown': {
      return (
        <Markdown components={markdownComponents} style={{ overflow: 'auto' }}>
          {content}
        </Markdown>
      );
    }

    default: {
      return <HTMLRenderer htmlContent={content} />;
    }
  }
});

export default Renderer;
