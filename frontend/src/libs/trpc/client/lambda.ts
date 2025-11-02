// Mock implementations for tRPC client

// Mock observable implementation
const createMockObservable = <T>(subscriber: (observer: any) => void) => ({
  subscribe: (callbacks: any) => {
    const observer = {
      next: callbacks.next || (() => {}),
      error: callbacks.error || (() => {}),
      complete: callbacks.complete || (() => {}),
    };
    subscriber(observer);
    return { unsubscribe: () => {} };
  },
});

// Mock TRPCLink
const mockTRPCLink = () => ({
  op: null as any,
  next: (op: any) => createMockObservable(() => {}),
});

// Mock httpBatchLink
const mockHttpBatchLink = (config: any) => mockTRPCLink();

// Mock createTRPCClient
const mockCreateTRPCClient = <T>(config: any) => ({
  query: () => Promise.resolve({}),
  mutation: () => Promise.resolve({}),
});

// Mock React Query hooks
const createMockMutation = (path?: string[]) => ({
  mutateAsync: async (params?: any) => {
    // Mock successful response - customize based on the route
    await new Promise(resolve => setTimeout(resolve, 100));

    // Special handling for exporter.exportPdf
    if (path && path.includes('exporter') && path.includes('exportPdf')) {
      // Generate mock PDF data
      const mockPdfBase64 = btoa(
        '%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n2 0 obj\n<<\n/Type /Pages\n/Kids [3 0 R]\n/Count 1\n>>\nendobj\n3 0 obj\n<<\n/Type /Page\n/Parent 2 0 R\n/MediaBox [0 0 612 792]\n/Contents 4 0 R\n>>\nendobj\n4 0 obj\n<<\n/Length 44\n>>\nstream\nBT\n/F1 12 Tf\n100 700 Td\n(Mock PDF Content) Tj\nET\nendstream\nendobj\nxref\n0 5\n0000000000 65535 f\n0000000009 00000 n\n0000000058 00000 n\n0000000115 00000 n\n0000000200 00000 n\ntrailer\n<<\n/Size 5\n/Root 1 0 R\n>>\nstartxref\n284\n%%EOF'
      );

      const mockFilename = params?.title
        ? `${params.title.replace(/[^a-zA-Z0-9]/g, '_')}.pdf`
        : 'document.pdf';

      return {
        pdf: mockPdfBase64,
        filename: mockFilename
      };
    }

    // Special handling for chunk.getChunksByFileId
    if (path && path.includes('chunk') && path.includes('getChunksByFileId')) {
      return {
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
              }
            ],
            nextCursor: null
          }
        ]
      };
    }

    // Default success response
    return { success: true };
  },
  isPending: false,
  error: null,
});

const createMockInfiniteQuery = () => ({
  data: { pages: [] },
  isLoading: false,
  error: null,
});

const createMockQuery = () => ({
  data: null,
  isLoading: false,
  error: null,
});

// Mock React hooks creator
const mockCreateTRPCReact = <T>() => {
  const createProxy = (path: string[] = []): any => {
    return new Proxy(() => {}, {
      get: (target, prop) => {
        if (prop === 'useMutation') {
          return () => createMockMutation(path);
        }
        if (prop === 'useInfiniteQuery') {
          return createMockInfiniteQuery;
        }
        if (prop === 'useQuery') {
          return createMockQuery;
        }
        // Return nested proxy for deeper paths
        return createProxy([...path, prop as string]);
      },
      apply: (target, thisArg, args) => {
        // Handle direct function calls
        return createProxy(path);
      },
    });
  };

  const trpc = createProxy();

  // Add createClient method
  trpc.createClient = (config: any) => ({
    query: () => Promise.resolve({}),
    mutation: () => Promise.resolve({}),
  });

  return trpc;
};

// Export mock implementations
export const TRPCLink = mockTRPCLink;
export const createTRPCClient = mockCreateTRPCClient;
export const httpBatchLink = mockHttpBatchLink;
export const observable = createMockObservable;
export const createTRPCReact = mockCreateTRPCReact;

// Create mock lambda client and query objects
export const lambdaClient = mockCreateTRPCClient<any>({
  links: [],
});

export const lambdaQuery = createTRPCReact<any>();

export const lambdaQueryClient = lambdaQuery.createClient({ links: [] });
