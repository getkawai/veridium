'use client';

import type { PDFDocumentProxy } from 'pdfjs-dist';
import { Fragment, memo, useCallback, useState } from 'react';
import { Flexbox } from 'react-layout-kit';
import { Document, Page, pdfjs } from 'react-pdf';
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';


import HighlightLayer from './HighlightLayer';
import { useStyles } from './style';
import useResizeObserver from './useResizeObserver';

// 如果海外的地址： https://unpkg.com/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.mjs
pdfjs.GlobalWorkerOptions.workerSrc = `https://registry.npmmirror.com/pdfjs-dist/${pdfjs.version}/files/build/pdf.worker.min.mjs`;

const options = {
  cMapUrl: `https://registry.npmmirror.com/pdfjs-dist/${pdfjs.version}/files/cmaps/`,
  standardFontDataUrl: `https://registry.npmmirror.com/pdfjs-dist/${pdfjs.version}/files/standard_fonts/`,
};

const maxWidth = 1200;

interface PDFViewerProps {
  fileId: string;
  url: string | null;
}

const PDFViewer = memo<PDFViewerProps>(({ url, fileId }) => {
  const { styles } = useStyles();
  const [numPages, setNumPages] = useState<number>(0);
  const [containerRef, setContainerRef] = useState<HTMLElement | null>(null);
  const [containerWidth, setContainerWidth] = useState<number>();
  const [isLoaded, setIsLoaded] = useState(false);

  // eslint-disable-next-line no-undef
  const onResize = useCallback<ResizeObserverCallback>((entries) => {
    const [entry] = entries;

    if (entry) {
      setContainerWidth(entry.contentRect.width);
    }
  }, []);

  useResizeObserver(containerRef, onResize);

  const onDocumentLoadSuccess = ({ numPages: nextNumPages }: PDFDocumentProxy) => {
    setNumPages(nextNumPages);
    setIsLoaded(true);
  };

  // Mock data for PDF chunks
  const mockData = {
    pages: [
      {
        items: [
          {
            id: 'chunk-1',
            text: 'This is the first paragraph of content that would be highlighted in the PDF.',
            pageNumber: 1,
            metadata: {
              coordinates: {
                layout_height: 842,
                layout_width: 595,
                points: [[50, 100], [545, 100], [545, 150], [50, 150]],
                system: 'pdf'
              },
              languages: ['en'],
              pageNumber: 1,
              text_as_html: '<p>This is the first paragraph of content that would be highlighted in the PDF.</p>'
            },
            createdAt: new Date('2024-01-01T00:00:00Z'),
            updatedAt: new Date('2024-01-01T00:00:00Z'),
            index: 0,
            type: 'text',
            parentId: null
          },
          {
            id: 'chunk-2',
            text: 'Here is another section with important information that should be highlighted.',
            pageNumber: 1,
            metadata: {
              coordinates: {
                layout_height: 842,
                layout_width: 595,
                points: [[50, 200], [545, 200], [545, 250], [50, 250]],
                system: 'pdf'
              },
              languages: ['en'],
              pageNumber: 1,
              text_as_html: '<p>Here is another section with important information that should be highlighted.</p>'
            },
            createdAt: new Date('2024-01-01T00:00:00Z'),
            updatedAt: new Date('2024-01-01T00:00:00Z'),
            index: 1,
            type: 'text',
            parentId: null
          }
        ],
        nextCursor: null
      },
      {
        items: [
          {
            id: 'chunk-3',
            text: 'Content on the second page that demonstrates multi-page highlighting capabilities.',
            pageNumber: 2,
            metadata: {
              coordinates: {
                layout_height: 842,
                layout_width: 595,
                points: [[50, 100], [545, 100], [545, 180], [50, 180]],
                system: 'pdf'
              },
              languages: ['en'],
              pageNumber: 2,
              text_as_html: '<p>Content on the second page that demonstrates multi-page highlighting capabilities.</p>'
            },
            createdAt: new Date('2024-01-01T00:00:00Z'),
            updatedAt: new Date('2024-01-01T00:00:00Z'),
            index: 2,
            type: 'text',
            parentId: null
          }
        ],
        nextCursor: null
      }
    ]
  };

  const dataSource = mockData.pages.flatMap((page) => page.items);

  return (
    <Flexbox className={styles.container}>
      <Flexbox
        align={'center'}
        className={styles.documentContainer}
        padding={24}
        ref={setContainerRef}
        style={{ height: isLoaded ? undefined : '100%' }}
      >
        <Document
          className={styles.document}
          file={url}
          onLoadSuccess={onDocumentLoadSuccess}
          options={options}
        >
          {Array.from({ length: numPages }, (el, index) => {
            const width = containerWidth ? Math.min(containerWidth, maxWidth) : maxWidth;

            return (
              <Fragment key={`page_${index + 1}`}>
                <Page className={styles.page} pageNumber={index + 1} width={width}>
                  <HighlightLayer dataSource={dataSource} pageNumber={index + 1} width={width} />
                </Page>
              </Fragment>
            );
          })}
        </Document>
      </Flexbox>
    </Flexbox>
  );
});

export default PDFViewer;
