import { BuiltinToolManifest } from '@/types';

export const VideoDescribeApiName = {
  getVideoTranscription: 'getVideoTranscription',
};

export const VideoDescribeManifest: BuiltinToolManifest = {
  api: [
    {
      description:
        'Get AI-generated transcription of an uploaded video audio. Use this when user asks about what is said in the video, spoken words, dialogue, or audio content. The transcription is generated using Whisper STT when the video was uploaded.',
      name: VideoDescribeApiName.getVideoTranscription,
      parameters: {
        properties: {
          file_id: {
            description: 'The file ID of the uploaded video',
            type: 'string',
          },
        },
        required: ['file_id'],
        type: 'object',
      },
    },
  ],
  identifier: 'lobe-video-describe',
  meta: {
    avatar: '🎬',
    description: 'Get AI-generated transcriptions of uploaded video audio',
    title: 'Video Describe',
  },
  systemRole: `You have access to the Video Describe tool which retrieves AI-generated transcriptions of uploaded video audio.

When a user uploads a video and asks about what is said or spoken in it, use the getVideoTranscription tool with the file_id to get the pre-generated transcription.

The transcription includes:
- Spoken words and dialogue from the video audio
- Automatic speech recognition using Whisper STT

Use this information to answer user questions about the video audio content accurately.`,
  type: 'builtin',
};
