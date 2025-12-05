import { BuiltinRenderProps } from '@/types';
import { Markdown } from '@lobehub/ui';
import { useTheme } from 'antd-style';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

interface ImageDescribeResponse {
  file_id: string;
  description: string;
  status: string;
}

interface ImageDescribeParams {
  file_id: string;
}

const ImageDescribeRender = memo<BuiltinRenderProps<ImageDescribeResponse, ImageDescribeParams>>(
  ({ content }) => {
    const theme = useTheme();

    if (!content) {
      return null;
    }

    return (
      <Flexbox gap={8}>
        <div
          style={{
            backgroundColor: theme.colorBgContainer,
            border: `1px solid ${theme.colorBorder}`,
            borderRadius: theme.borderRadius,
            maxHeight: 400,
            overflow: 'auto',
            padding: 12,
          }}
        >
          <Markdown>{content.description || 'No description available'}</Markdown>
        </div>
      </Flexbox>
    );
  },
);

export default ImageDescribeRender;
