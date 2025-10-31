import { useCallback, useState } from 'react';

interface PdfGenerationParams {
  content: string;
  sessionId: string;
  title?: string;
  topicId?: string;
}

interface PdfGenerationState {
  downloadPdf: () => Promise<void>;
  error: string | null;
  generatePdf: (params: PdfGenerationParams) => Promise<void>;
  loading: boolean;
  pdfData: string | null;
}

export const usePdfGeneration = (): PdfGenerationState => {
  const [pdfData, setPdfData] = useState<string | null>(null);
  const [filename, setFilename] = useState<string>('chat-export.pdf');
  const [error, setError] = useState<string | null>(null);
  const [lastGeneratedKey, setLastGeneratedKey] = useState<string | null>(null);
  const [isGenerating, setIsGenerating] = useState<boolean>(false);

  // Mock PDF generation - creates a minimal valid PDF as base64
  const mockExportPdfMutation = {
    isPending: isGenerating,
    error: error ? { message: error } : null,
    mutateAsync: async (params: PdfGenerationParams) => {
      setIsGenerating(true);
      setError(null);

      // Simulate network delay
      await new Promise(resolve => setTimeout(resolve, 1500));

      // Create mock PDF filename
      const title = params.title || 'Chat Export';
      const sessionId = params.sessionId.slice(-8); // Use last 8 chars of session ID
      const mockFilename = `${title.replace(/[^a-zA-Z0-9]/g, '_')}_${sessionId}.pdf`;

      // Mock PDF data - minimal valid PDF as base64
      // This is a very basic PDF structure (not readable, but valid format)
      const mockPdfData = btoa(
        '%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n2 0 obj\n<<\n/Type /Pages\n/Kids [3 0 R]\n/Count 1\n>>\nendobj\n3 0 obj\n<<\n/Type /Page\n/Parent 2 0 R\n/MediaBox [0 0 612 792]\n/Contents 4 0 R\n>>\nendobj\n4 0 obj\n<<\n/Length 44\n>>\nstream\nBT\n/F1 12 Tf\n100 700 Td\n(Mock PDF Content) Tj\nET\nendstream\nendobj\nxref\n0 5\n0000000000 65535 f\n0000000009 00000 n\n0000000058 00000 n\n0000000115 00000 n\n0000000200 00000 n\ntrailer\n<<\n/Size 5\n/Root 1 0 R\n>>\nstartxref\n284\n%%EOF'
      );

      setIsGenerating(false);
      return {
        pdf: mockPdfData,
        filename: mockFilename,
      };
    },
  };

  const generatePdf = useCallback(
    async (params: PdfGenerationParams) => {
      const { content, sessionId, title, topicId } = params;
      // Create a key to identify this specific request
      const requestKey = `${sessionId}-${topicId || 'default'}-${content.length}`;

      // Prevent multiple simultaneous requests or re-generating the same PDF
      if (mockExportPdfMutation.isPending || lastGeneratedKey === requestKey) return;

      try {
        setError(null);
        setPdfData(null);

        const result = await mockExportPdfMutation.mutateAsync({
          content,
          sessionId,
          title,
          topicId,
        });

        setPdfData(result.pdf);
        setFilename(result.filename);
        setLastGeneratedKey(requestKey);
      } catch (error) {
        console.error('Failed to generate PDF:', error);
        setError(error instanceof Error ? error.message : 'Failed to generate PDF');
      }
    },
    [lastGeneratedKey],
  );

  const downloadPdf = useCallback(async () => {
    if (!pdfData) return;

    try {
      // Convert base64 to blob
      const byteCharacters = atob(pdfData);
      const byteNumbers = Array.from({ length: byteCharacters.length }, (_, i) =>
        byteCharacters.charCodeAt(i),
      );
      const byteArray = new Uint8Array(byteNumbers);
      const blob = new Blob([byteArray], { type: 'application/pdf' });

      // Create download link
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.append(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to download PDF:', error);
      throw error;
    }
  }, [pdfData, filename]);

  return {
    downloadPdf,
    error: error || (mockExportPdfMutation.error?.message ?? null),
    generatePdf,
    loading: mockExportPdfMutation.isPending,
    pdfData,
  };
};
