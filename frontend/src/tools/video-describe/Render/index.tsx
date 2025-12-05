import { BuiltinRenderProps } from '@/types';
import { Markdown } from '@lobehub/ui';
import { useTheme } from 'antd-style';
import { memo } from 'react';
import { Flexbox } from 'react-layout-kit';

interface VideoDescribeResponse {
  file_id: string;
  transcription: string;
  status: string;
}

interface VideoDescribeParams {
  file_id: string;
}

const VideoDescribeRender = memo<BuiltinRenderProps<VideoDescribeResponse, VideoDescribeParams>>(
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
          <Markdown>{content.transcription || 'No transcription available'}</Markdown>
        </div>
      </Flexbox>
    );
  },
);

export default VideoDescribeRender;
