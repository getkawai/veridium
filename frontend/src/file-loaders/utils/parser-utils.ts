import { ExtractFiles as WailsExtractFiles } from '@@/github.com/kawai-network/veridium/zipservice';

// Define basic error messages
const ERRORMSG = {
  fileCorrupted: (filepath: string) =>
    `[OfficeParser]: Your file ${typeof filepath === 'string' ? filepath : 'Buffer'} seems to be corrupted. If you are sure it is fine, please create a ticket.`,
  invalidInput: `[OfficeParser]: Invalid input type: Expected a Buffer or a valid file path`,
};

/** Returns parsed xml document for a given xml text.
 * @param {string} xml The xml string from the doc file
 * @returns {XMLDocument}
 */
export const parseString = (xml: string) => {
  const parser = new window.DOMParser();
  return parser.parseFromString(xml, 'text/xml') as unknown as XMLDocument;
};

export interface ExtractedFile {
  content: string;
  path: string;
}

/** Extract specific files from a ZIP file based on a filter function.
 * @param {string} filePath ZIP file path (string).
 * @param {(fileName: string) => boolean} filterFn A function that receives the entry file name and returns true if the file should be extracted.
 * @returns {Promise<ExtractedFile[]>} Resolves to an array of object containing file path and content.
 */
export async function extractFiles(
  filePath: string,
  filterFn: (fileName: string) => boolean,
): Promise<ExtractedFile[]> {
  // For now, we only support string paths
  if (typeof filePath !== 'string') {
    throw new Error(ERRORMSG.invalidInput);
  }

  // Convert the filter function to a regex pattern
  // This is a heuristic approach since we can't execute JS functions in Go
  let pattern = '.*'; // Default: match everything

  const filterString = filterFn.toString();

  // Try to detect common patterns based on the filter function string
  if (filterString.includes('ppt/slides/slide') && filterString.includes('.xml')) {
    // PPTX slide pattern
    pattern = 'ppt/slides/slide\\d+\\.xml';
  } else if (filterString.includes('xl/worksheets/sheet') && filterString.includes('.xml')) {
    // Excel worksheet pattern
    pattern = 'xl/worksheets/sheet\\d+\\.xml';
  } else if (filterString.includes('word/document') && filterString.includes('.xml')) {
    // Word document pattern
    pattern = 'word/document\\d*\\.xml';
  }

  try {
    // Call the Wails service
    const result = await WailsExtractFiles(filePath, pattern);

    // Convert base64 content back to UTF-8 strings
    return result.map(file => ({
      content: atob(file.content), // Decode base64 to string
      path: file.path,
    }));
  } catch (error) {
    throw new Error(`Failed to extract files from ZIP: ${error}`);
  }
}
