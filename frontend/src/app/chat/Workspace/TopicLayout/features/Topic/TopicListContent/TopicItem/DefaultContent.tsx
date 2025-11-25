import { Icon, Tag, Text } from '@lobehub/ui';
import { useTheme } from 'antd-style';
import { MessageSquareDashed } from 'lucide-react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Flexbox } from 'react-layout-kit';

// MASIH DIGUNAKAN: Component untuk menampilkan UI dari topic "default" (temporary chat)
// - Dipanggil dari TopicItem/index.tsx ketika id === undefined (line 75: {!id ? <DefaultContent /> : ...})
// - Menampilkan icon MessageSquareDashed + text "默认话题" (defaultTitle) + tag "Temp"
// - Digunakan untuk menandai bahwa chat saat ini belum disimpan ke topic permanen
// - User bisa klik tombol "Save Topic" untuk mengkonversi temporary chat ini menjadi topic permanen
const DefaultContent = memo(() => {
  const { t } = useTranslation('topic');

  const theme = useTheme();

  return (
    <Flexbox align={'center'} gap={8} horizontal>
      <Flexbox align={'center'} height={24} justify={'center'} width={24}>
        <Icon color={theme.colorTextDescription} icon={MessageSquareDashed} />
      </Flexbox>
      <Text ellipsis={{ rows: 1 }} style={{ margin: 0 }}>
        {t('defaultTitle')}
      </Text>
      <Tag>{t('temp')}</Tag>
    </Flexbox>
  );
});

export default DefaultContent;
