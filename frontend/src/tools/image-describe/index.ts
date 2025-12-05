import { BuiltinToolManifest } from '@/types';

export const ImageDescribeApiName = {
  getImageDescription: 'getImageDescription',
};

export const ImageDescribeManifest: BuiltinToolManifest = {
  api: [
    {
      description:
        'Get AI-generated description of an uploaded image or video. Use this when user asks about image content, text extraction, OCR, or visual analysis. The description is pre-generated when the file was uploaded using local VL model.',
      name: ImageDescribeApiName.getImageDescription,
      parameters: {
        properties: {
          file_id: {
            description: 'The file ID of the uploaded image or video',
            type: 'string',
          },
        },
        required: ['file_id'],
        type: 'object',
      },
    },
  ],
  identifier: 'lobe-image-describe',
  meta: {
    avatar: '🖼️',
    description: 'Get AI-generated descriptions of uploaded images and videos',
    title: 'Image Describe',
  },
  systemRole: `You have access to the Image Describe tool which retrieves AI-generated descriptions of uploaded images and videos.

When a user uploads an image or video and asks about its content, use the getImageDescription tool with the file_id to get the pre-generated description.

The description includes:
- Visual content analysis
- Text extraction (OCR) if text is present
- Object and scene recognition
- Layout and structure information

Use this information to answer user questions about the uploaded files accurately.`,
  type: 'builtin',
};
