import { BuiltinRender } from '@/types';

import { CodeInterpreterManifest } from './code-interpreter';
import CodeInterpreterRender from './code-interpreter/Render';
import { DalleManifest } from './dalle';
import DalleRender from './dalle/Render';
import { ImageDescribeManifest } from './image-describe';
import ImageDescribeRender from './image-describe/Render';
import { LocalSystemManifest } from './local-system';
import LocalFilesRender from './local-system/Render';
import { VideoDescribeManifest } from './video-describe';
import VideoDescribeRender from './video-describe/Render';
import { WebBrowsingManifest } from './web-browsing';
import WebBrowsing from './web-browsing/Render';

export const BuiltinToolsRenders: Record<string, BuiltinRender> = {
  [DalleManifest.identifier]: DalleRender as BuiltinRender,
  [WebBrowsingManifest.identifier]: WebBrowsing as BuiltinRender,
  [LocalSystemManifest.identifier]: LocalFilesRender as BuiltinRender,
  [CodeInterpreterManifest.identifier]: CodeInterpreterRender as BuiltinRender,
  [ImageDescribeManifest.identifier]: ImageDescribeRender as BuiltinRender,
  [VideoDescribeManifest.identifier]: VideoDescribeRender as BuiltinRender,
};
