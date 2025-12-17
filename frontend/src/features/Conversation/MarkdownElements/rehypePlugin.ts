
import { SKIP, visit } from 'unist-util-visit';

export const createRehypePlugin = (tagName: string) => {
  return () => {
    return (tree: any) => {
      visit(tree, (node, index, parent) => {
        if (node.type === 'element' && node.tagName === 'p' && node.children.length > 0) {
          // Strategy 1: Check for existing Element node
          const elementNodeIndex = node.children.findIndex((child: any) =>
            child.type === 'element' && child.tagName === tagName
          );

          if (elementNodeIndex !== -1) {
            const elementNode = node.children[elementNodeIndex];
            // Check if others are whitespace
            const isOnlyContent = node.children.every((child: any, i: number) => {
              if (i === elementNodeIndex) return true;
              return (child.type === 'text' && !child.value.trim()) ||
                (child.type === 'raw' && !child.value.trim());
            });

            if (isOnlyContent) {
              parent.children.splice(index, 1, elementNode);
              return [SKIP, index];
            }
          }

          // Strategy 2: Check for Raw/Text nodes forming a wrapper block (even if split)
          const textContent = node.children.map((c: any) => c.value || '').join('').trim();

          // Regex to match the full block pattern: ^<tagName ...> ... </tagName>$
          // We use [\s\S]*? for non-greedy content match across lines
          const wrapperRegex = new RegExp(`^<${tagName}\\b[^>]*>([\\s\\S]*?)<\\/${tagName}>$`);

          if (wrapperRegex.test(textContent)) {
            // Extract attributes from the opening tag part
            const attributes: Record<string, string> = {};
            const attributeRegex = /(\w+)="([^"]*)"/g;
            let match;

            const openingTagMatch = textContent.match(new RegExp(`^<${tagName}\\b[^>]*>`));
            if (openingTagMatch) {
              const openingTag = openingTagMatch[0];
              while ((match = attributeRegex.exec(openingTag)) !== null) {
                attributes[match[1]] = match[2];
              }
            }

            // Extract inner content matching the regex group
            const innerMatch = textContent.match(wrapperRegex);
            const innerContent = innerMatch ? innerMatch[1] : '';

            // Construct new element node
            const newNode = {
              children: [
                {
                  type: 'text',
                  value: innerContent,
                }
              ],
              properties: attributes,
              tagName: tagName,
              type: 'element',
            };

            parent.children.splice(index, 1, newNode);
            return [SKIP, index];
          }
        }
        else if (node.type === 'raw' && node.value.trim().startsWith(`<${tagName}`)) {
          // Handle raw node case (top level raw block, or inside other elements but not caught by p check)
          // We need to parse attributes here too if possible, but raw nodes usually don't expose attributes easily 
          // unless we parse the string.
          // However, for hydration purposes, just converting to element is often enough.

          const content = node.value.trim();
          // Try to parse attributes
          const attributes: Record<string, string> = {};
          const attributeRegex = /(\w+)="([^"]*)"/g;
          let match;

          // Match opening tag
          const openingTagMatch = content.match(new RegExp(`^<${tagName}\\b[^>]*>`));
          if (openingTagMatch) {
            const openingTag = openingTagMatch[0];
            while ((match = attributeRegex.exec(openingTag)) !== null) {
              attributes[match[1]] = match[2];
            }
          }

          const newNode = {
            // We capture the whole raw content as a text child for the new element.
            // The component renderer (React) will likely parse this content or use it.
            children: [{
              type: 'text',
              value: content
            }],
            properties: attributes,
            tagName: tagName,
            type: 'element',
          };

          // Actually, parsing content from a single raw text node:
          const wrapperRegex = new RegExp(`^<${tagName}\\b[^>]*>([\\s\\S]*?)<\\/${tagName}>$`);
          const innerMatch = content.match(wrapperRegex);
          if (innerMatch) {
            newNode.children = [{ type: 'text', value: innerMatch[1] }];
          }

          parent.children.splice(index, 1, newNode);
          return [SKIP, index];
        }
      });
    };
  };
};
