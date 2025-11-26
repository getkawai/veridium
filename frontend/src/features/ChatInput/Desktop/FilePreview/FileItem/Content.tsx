import { Image } from '@lobehub/ui';
import { createStyles } from 'antd-style';
import { memo } from 'react';

import FileIcon from '@/components/FileIcon';
import { UploadFileItem } from '@/types/files/upload';

const useStyles = createStyles(({ css }) => ({
  image: css`
    margin-block: 0 !important;
    box-shadow: none;

    img {
      object-fit: contain;
    }
  `,
  video: css`
    overflow: hidden;
    border-radius: 8px;
  `,
}));

const Content = memo<UploadFileItem>(({ file, previewUrl, base64Url }) => {
  const { styles } = useStyles();

  // Use base64Url for Wails desktop app (Blob URLs don't work)
  // Fall back to previewUrl for web
  const imageSource = base64Url || previewUrl;
  
  console.log('[Content] Rendering file:', file.name, {
    fileType: file.type,
    hasPreviewUrl: !!previewUrl,
    hasBase64Url: !!base64Url,
    imageSource: imageSource?.substring(0, 50) + '...',
  });

  if (file.type.startsWith('image')) {
    return <Image alt={file.name} src={imageSource} wrapperClassName={styles.image} />;
  }

  if (file.type.startsWith('video')) {
    return <video className={styles.video} src={imageSource} width={'100%'} />;
  }

  return <FileIcon fileName={file.name} fileType={file.type} size={48} />;
});

export default Content;
