# Artifacts Plugin/Tool Documentation

## Overview

Artifacts is a built-in plugin/tool in LobeChat that enables AI assistants to create substantial, self-contained content that users can view, modify, and reuse. Inspired by Claude's Artifacts feature, it provides a separate UI window for displaying rich content like code, HTML pages, React components, SVG graphics, diagrams, and documents.

## Purpose

The Artifacts plugin serves to:

- **Separate substantial content** from conversational flow for better clarity
- **Enable live rendering** of code, HTML, React components, and graphics
- **Facilitate iteration** on user-generated content
- **Provide interactive previews** of code and visualizations
- **Support multiple content types** with appropriate renderers

## Key Concepts

### What Makes Good Artifacts?

✅ **Good artifacts are:**
- Substantial content (>15 lines)
- Content users might modify, iterate on, or take ownership of
- Self-contained, complex content understandable on its own
- Content intended for eventual use outside the conversation
- Content likely to be referenced or reused multiple times

❌ **Don't use artifacts for:**
- Simple, informational, or short content
- Brief code snippets or small examples
- Explanatory or instructional content
- Suggestions or feedback on existing artifacts
- Conversational content dependent on context
- Content unlikely to be modified
- One-off questions

## Architecture

### Component Hierarchy

```
Artifacts System
├── Tool Definition (src/tools/artifacts/)
│   ├── index.ts - Manifest definition
│   └── systemRole.ts - AI instructions
│
├── Markdown Parser (src/features/Conversation/MarkdownElements/LobeArtifact/)
│   ├── rehypePlugin.ts - Parse <lobeArtifact> tags
│   ├── Render/index.tsx - Inline artifact card
│   └── Render/Icon.tsx - Type-specific icons
│
├── Portal Display (src/features/Portal/Artifacts/)
│   ├── Body/index.tsx - Main artifact viewer
│   ├── Body/Renderer/index.tsx - Type router
│   ├── Body/Renderer/React/index.tsx - React renderer
│   ├── Body/Renderer/HTML.tsx - HTML renderer
│   ├── Body/Renderer/SVG.tsx - SVG renderer
│   └── Header.tsx - Artifact header
│
└── State Management (src/store/chat/slices/portal/)
    ├── action.ts - Portal actions
    ├── selectors.ts - Portal selectors
    └── initialState.ts - Portal state
```

### Data Flow

```
AI generates message with <lobeArtifact> tag
                ↓
Markdown parser detects artifact tag (rehypePlugin)
                ↓
Extracts attributes (identifier, type, title, language)
                ↓
Renders inline artifact card in chat
                ↓
User clicks artifact card
                ↓
openArtifact() action called
                ↓
Portal opens with artifact viewer
                ↓
Renderer component selected based on type
                ↓
Content rendered (React/HTML/SVG/Mermaid/Markdown/Code)
```

## Artifact Types

### 1. Code (`application/lobe.artifacts.code`)

For code snippets or scripts in any programming language.

**Attributes**:
- `type`: `"application/lobe.artifacts.code"`
- `language`: Programming language (e.g., `"python"`, `"javascript"`)
- `identifier`: Unique kebab-case identifier
- `title`: Brief description

**Example**:

```xml
<lobeArtifact identifier="factorial-script" type="application/lobe.artifacts.code" language="python" title="Factorial Calculator">
def factorial(n):
    if n == 0:
        return 1
    else:
        return n * factorial(n - 1)

number = int(input("Enter a number: "))
print(f"Factorial of {number} is {factorial(number)}")
</lobeArtifact>
```

**Rendering**: Displays syntax-highlighted code using `Highlighter` component.

### 2. HTML (`text/html`)

For single-file HTML pages with embedded CSS and JavaScript.

**Attributes**:
- `type`: `"text/html"`
- `identifier`: Unique identifier
- `title`: Page title

**Constraints**:
- Must be single-file (HTML + CSS + JS combined)
- External scripts only from `https://cdnjs.cloudflare.com`
- Images: Use placeholder API `/api/placeholder/{width}/{height}`
- No external image URLs allowed

**Example**:

```xml
<lobeArtifact identifier="interactive-calculator" type="text/html" title="Simple Calculator">
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: Arial; text-align: center; padding: 50px; }
    button { margin: 5px; padding: 10px 20px; font-size: 18px; }
  </style>
</head>
<body>
  <h1>Calculator</h1>
  <input id="display" type="text" readonly style="font-size: 24px; padding: 10px;">
  <div>
    <button onclick="appendNumber('1')">1</button>
    <button onclick="appendNumber('2')">2</button>
    <!-- more buttons -->
  </div>
  <script>
    function appendNumber(num) {
      document.getElementById('display').value += num;
    }
  </script>
</body>
</html>
</lobeArtifact>
```

**Rendering**: Uses iframe with `HTMLRenderer` component to safely render HTML.

### 3. SVG (`image/svg+xml`)

For Scalable Vector Graphics images.

**Attributes**:
- `type`: `"image/svg+xml"`
- `identifier`: Unique identifier
- `title`: Image description

**Best Practices**:
- Use `viewBox` instead of fixed `width`/`height`
- Keep SVG self-contained
- Use semantic element IDs

**Example**:

```xml
<lobeArtifact identifier="blue-circle" type="image/svg+xml" title="Simple Blue Circle">
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <circle cx="50" cy="50" r="40" fill="blue" />
</svg>
</lobeArtifact>
```

**Rendering**: Direct SVG rendering with interactive features via `SVGRender` component.

### 4. Mermaid Diagrams (`application/lobe.artifacts.mermaid`)

For flowcharts, sequence diagrams, and other Mermaid visualizations.

**Attributes**:
- `type`: `"application/lobe.artifacts.mermaid"`
- `identifier`: Unique identifier
- `title`: Diagram description

**Example**:

```xml
<lobeArtifact identifier="user-flow" type="application/lobe.artifacts.mermaid" title="User Login Flow">
graph TD
    A[Start] --> B{Logged in?}
    B -->|Yes| C[Dashboard]
    B -->|No| D[Login Page]
    D --> E{Valid credentials?}
    E -->|Yes| C
    E -->|No| D
    C --> F[End]
</lobeArtifact>
```

**Rendering**: Uses `Mermaid` component from `@lobehub/ui`.

### 5. Markdown Documents (`text/markdown`)

For formatted text documents.

**Attributes**:
- `type`: `"text/markdown"`
- `identifier`: Unique identifier
- `title`: Document title

**Example**:

```xml
<lobeArtifact identifier="project-proposal" type="text/markdown" title="Project Proposal">
# Project Proposal

## Executive Summary
This proposal outlines...

## Objectives
1. Increase efficiency
2. Reduce costs
3. Improve user experience

## Timeline
- Q1: Planning
- Q2: Implementation
- Q3: Testing
- Q4: Launch
</lobeArtifact>
```

**Rendering**: Uses `Markdown` component with overflow handling.

### 6. React Components (`application/lobe.artifacts.react`)

For interactive React components with live preview.

**Attributes**:
- `type`: `"application/lobe.artifacts.react"`
- `identifier`: Unique identifier
- `title`: Component name

**Available Libraries**:
- `react` - Base React with hooks
- `lucide-react@0.263.1` - Icon library
- `recharts` - Charting library
- `@/components/ui/*` - shadcn/ui components (Alert, Button, Card, etc.)
- `tailwindcss` - Styling (via CDN)

**Constraints**:
- Component must have no required props or provide defaults
- Must use default export
- Use Tailwind classes for styling (NO arbitrary values like `h-[600px]`)
- No other libraries (zod, hookform, etc.) available
- Images: Use placeholder API only

**Example**:

```xml
<lobeArtifact identifier="todo-app" type="application/lobe.artifacts.react" title="Todo List App">
import React, { useState } from 'react';
import { Plus, Trash2 } from 'lucide-react';

export default function TodoApp() {
  const [todos, setTodos] = useState([]);
  const [input, setInput] = useState('');

  const addTodo = () => {
    if (input.trim()) {
      setTodos([...todos, { id: Date.now(), text: input, done: false }]);
      setInput('');
    }
  };

  const toggleTodo = (id) => {
    setTodos(todos.map(t => t.id === id ? { ...t, done: !t.done } : t));
  };

  const deleteTodo = (id) => {
    setTodos(todos.filter(t => t.id !== id));
  };

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg">
      <h1 className="text-2xl font-bold mb-4">My Todos</h1>
      <div className="flex gap-2 mb-4">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && addTodo()}
          className="flex-1 px-3 py-2 border rounded"
          placeholder="Add a todo..."
        />
        <button onClick={addTodo} className="px-4 py-2 bg-blue-500 text-white rounded">
          <Plus size={20} />
        </button>
      </div>
      <ul className="space-y-2">
        {todos.map(todo => (
          <li key={todo.id} className="flex items-center gap-2 p-2 border rounded">
            <input
              type="checkbox"
              checked={todo.done}
              onChange={() => toggleTodo(todo.id)}
            />
            <span className={todo.done ? 'line-through flex-1' : 'flex-1'}>
              {todo.text}
            </span>
            <button onClick={() => deleteTodo(todo.id)} className="text-red-500">
              <Trash2 size={16} />
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
</lobeArtifact>
```

**Rendering**: Uses Sandpack (CodeSandbox) with Vite + React + TypeScript template.

## Implementation Details

### 1. Tool Manifest

**File**: `src/tools/artifacts/index.ts`

```typescript
export const ArtifactsManifest: BuiltinToolManifest = {
  api: [],
  identifier: 'lobe-artifacts',
  meta: {
    avatar: `data:image/svg+xml;base64,...`,
    title: 'Artifacts',
  },
  systemRole: systemPrompt,
  type: 'builtin',
};
```

**Key Points**:
- Registered as a built-in tool
- No API endpoints (client-side only)
- Includes comprehensive system prompt for AI

### 2. System Role (AI Instructions)

**File**: `src/tools/artifacts/systemRole.ts`

The system prompt provides detailed instructions to the AI on:
- When to create artifacts vs. inline content
- How to structure artifact tags
- Available artifact types and their constraints
- Update vs. create logic (reuse identifiers)
- Multiple examples of correct usage

**Key Instruction Pattern**:

```xml
<lobeThinking>
Evaluate if content is artifact-worthy. 
If updating, reuse identifier. If new, create unique identifier.
</lobeThinking>

<lobeArtifact identifier="unique-id" type="type" title="Title">
Content here
</lobeArtifact>
```

### 3. Markdown Parser (Rehype Plugin)

**File**: `src/features/Conversation/MarkdownElements/LobeArtifact/rehypePlugin.ts`

Parses `<lobeArtifact>` tags from markdown and converts them to React components.

**Process**:

```typescript
function rehypeAntArtifact() {
  return (tree: any) => {
    visit(tree, (node, index, parent) => {
      // Detect <lobeArtifact> tag in paragraph
      if (node.type === 'element' && node.tagName === 'p') {
        const firstChild = node.children[0];
        if (firstChild.type === 'raw' && firstChild.value.startsWith('<lobeArtifact')) {
          // Extract attributes using regex
          const attributes = extractAttributes(firstChild.value);
          
          // Extract content between tags
          const content = extractContent(node.children);
          
          // Create new artifact node
          const newNode = {
            type: 'element',
            tagName: 'lobeArtifact',
            properties: attributes,
            children: [{ type: 'text', value: content }],
          };
          
          // Replace paragraph with artifact node
          parent.children.splice(index, 1, newNode);
        }
      }
    });
  };
}
```

**Regex Patterns** (from `packages/const/src/plugin.ts`):

```typescript
// Match artifact tag with content
export const ARTIFACT_TAG_REGEX = 
  /<lobeArtifact\b[^>]*>(?<content>[\S\s]*?)(?:<\/lobeArtifact>|$)/;

// Check if artifact tag is closed
export const ARTIFACT_TAG_CLOSED_REGEX = 
  /<lobeArtifact\b[^>]*>([\S\s]*?)<\/lobeArtifact>/;
```

### 4. Inline Artifact Card

**File**: `src/features/Conversation/MarkdownElements/LobeArtifact/Render/index.tsx`

Displays a clickable card in the chat for each artifact.

**Features**:
- Shows artifact icon based on type
- Displays title and identifier
- Shows content length and loading state
- Click to open in portal
- Auto-opens when generating

**Component Structure**:

```tsx
const Render = memo<ArtifactProps>(({ identifier, title, type, language, children, id }) => {
  const [isGenerating, isArtifactTagClosed, openArtifact, closeArtifact] = useChatStore(...);

  // Auto-open artifact when generating
  useEffect(() => {
    if (!hasChildren || !isGenerating) return;
    openArtifact({ id, identifier, language, title, type });
  }, [isGenerating, hasChildren, ...]);

  return (
    <Flexbox
      className={styles.container}
      onClick={() => {
        const currentId = chatPortalSelectors.artifactMessageId(useChatStore.getState());
        if (currentId === id) {
          closeArtifact(); // Toggle off if already open
        } else {
          openArtifact({ id, identifier, language, title, type });
        }
      }}
    >
      <Center className={styles.avatar}>
        <ArtifactIcon type={type} />
      </Center>
      <Flexbox>
        <div>{title || 'Generating...'}</div>
        <div>{identifier} · {contentLength}</div>
      </Flexbox>
    </Flexbox>
  );
});
```

### 5. Portal State Management

**File**: `src/store/chat/slices/portal/action.ts`

Manages artifact display state.

**Actions**:

```typescript
export interface ChatPortalAction {
  openArtifact: (artifact: PortalArtifact) => void;
  closeArtifact: () => void;
  togglePortal: (open?: boolean) => void;
}

export const chatPortalSlice: StateCreator<ChatStore, [], [], ChatPortalAction> = 
  (set, get) => ({
    openArtifact: (artifact) => {
      get().togglePortal(true);
      set({ portalArtifact: artifact }, false, 'openArtifact');
    },
    
    closeArtifact: () => {
      get().togglePortal(false);
      set({ portalArtifact: undefined }, false, 'closeArtifact');
    },
    
    togglePortal: (open) => {
      const showInspector = open === undefined ? !get().showPortal : open;
      set({ showPortal: showInspector }, false, 'toggleInspector');
    },
  });
```

**State Interface**:

```typescript
export interface ChatPortalState {
  portalArtifact?: PortalArtifact;
  portalArtifactDisplayMode?: ArtifactDisplayMode; // 'code' | 'preview'
  showPortal: boolean;
}

export interface PortalArtifact {
  id: string;              // Message ID
  identifier?: string;     // Artifact identifier
  title?: string;          // Artifact title
  type?: string;           // MIME type
  language?: string;       // Programming language
  children?: string;       // Content
}
```

**Selectors** (`src/store/chat/slices/portal/selectors.ts`):

```typescript
export const chatPortalSelectors = {
  // Portal visibility
  showPortal: (s: ChatStoreState) => s.showPortal,
  showArtifactUI: (s: ChatStoreState) => !!s.portalArtifact,
  
  // Artifact metadata
  artifactTitle: (s: ChatStoreState) => s.portalArtifact?.title,
  artifactIdentifier: (s: ChatStoreState) => s.portalArtifact?.identifier || '',
  artifactMessageId: (s: ChatStoreState) => s.portalArtifact?.id,
  artifactType: (s: ChatStoreState) => s.portalArtifact?.type,
  artifactCodeLanguage: (s: ChatStoreState) => s.portalArtifact?.language,
  
  // Extract artifact content from message
  artifactCode: (id: string) => (s: ChatStoreState) => {
    const messageContent = chatSelectors.getMessageById(id)(s)?.content || '';
    const result = messageContent.match(ARTIFACT_TAG_REGEX);
    let content = result?.groups?.content || '';
    
    // Remove markdown code block wrapper if present
    content = content.replace(/^\s*```[^\n]*\n([\S\s]*?)\n```\s*$/, '$1');
    
    return content;
  },
  
  // Check if artifact tag is closed
  isArtifactTagClosed: (id: string) => (s: ChatStoreState) => {
    const content = chatSelectors.getMessageById(id)(s)?.content || '';
    return ARTIFACT_TAG_CLOSED_REGEX.test(content || '');
  },
};
```

### 6. Artifact Viewer (Portal)

**File**: `src/features/Portal/Artifacts/Body/index.tsx`

Main component that displays the artifact in the portal.

```typescript
const ArtifactsUI = () => {
  const [
    messageId,
    displayMode,
    isMessageGenerating,
    artifactType,
    artifactContent,
    artifactCodeLanguage,
    isArtifactTagClosed,
  ] = useChatStore((s) => {
    const messageId = chatPortalSelectors.artifactMessageId(s) || '';
    return [
      messageId,
      s.portalArtifactDisplayMode,
      chatSelectors.isMessageGenerating(messageId)(s),
      chatPortalSelectors.artifactType(s),
      chatPortalSelectors.artifactCode(messageId)(s),
      chatPortalSelectors.artifactCodeLanguage(s),
      chatPortalSelectors.isArtifactTagClosed(messageId)(s),
    ];
  });

  // Auto-switch to preview when artifact tag closes
  useEffect(() => {
    if (isMessageGenerating && isArtifactTagClosed && displayMode === 'code') {
      useChatStore.setState({ portalArtifactDisplayMode: 'preview' });
    }
  }, [isMessageGenerating, displayMode, isArtifactTagClosed]);

  // Determine language for syntax highlighting
  const language = useMemo(() => {
    switch (artifactType) {
      case 'application/lobe.artifacts.react': return 'tsx';
      case 'application/lobe.artifacts.code': return artifactCodeLanguage;
      case 'python': return 'python';
      default: return 'html';
    }
  }, [artifactType, artifactCodeLanguage]);

  // Show code view when:
  // - Artifact tag not closed yet
  // - Display mode is 'code'
  // - Artifact type is 'code' (non-renderable)
  const showCode = !isArtifactTagClosed || 
                   displayMode === 'code' || 
                   artifactType === 'application/lobe.artifacts.code';

  return (
    <Flexbox className="portal-artifact" flex={1} height="100%">
      {showCode ? (
        <Highlighter language={language} style={{ height: '100%' }}>
          {artifactContent}
        </Highlighter>
      ) : (
        <Renderer content={artifactContent} type={artifactType} />
      )}
    </Flexbox>
  );
};
```

### 7. Renderer Router

**File**: `src/features/Portal/Artifacts/Body/Renderer/index.tsx`

Routes to appropriate renderer based on artifact type.

```typescript
const Renderer = memo<{ content: string; type?: string }>(({ content, type }) => {
  switch (type) {
    case 'application/lobe.artifacts.react':
      return <ReactRenderer code={content} />;

    case 'image/svg+xml':
      return <SVGRender content={content} />;

    case 'application/lobe.artifacts.mermaid':
      return <Mermaid variant="borderless">{content}</Mermaid>;

    case 'text/markdown':
      return <Markdown style={{ overflow: 'auto' }}>{content}</Markdown>;

    default: // HTML and other types
      return <HTMLRenderer htmlContent={content} />;
  }
});
```

### 8. React Renderer (Sandpack)

**File**: `src/features/Portal/Artifacts/Body/Renderer/React/index.tsx`

Renders React components using CodeSandbox's Sandpack.

```typescript
const ReactRenderer = memo<{ code: string }>(({ code }) => {
  const title = useChatStore(chatPortalSelectors.artifactTitle);

  return (
    <SandpackProvider
      template="vite-react-ts"
      theme="auto"
      files={{
        'App.tsx': code,
        ...createTemplateFiles({ title }),
      }}
      customSetup={{
        dependencies: {
          'react': 'latest',
          'lucide-react': 'latest',
          'recharts': 'latest',
          '@lshay/ui': 'latest',
          'antd': 'latest',
          '@radix-ui/react-alert-dialog': 'latest',
          '@radix-ui/react-dialog': 'latest',
          '@radix-ui/react-icons': 'latest',
          'class-variance-authority': 'latest',
          'clsx': 'latest',
          'tailwind-merge': 'latest',
        },
      }}
      options={{
        externalResources: ['https://cdn.tailwindcss.com'],
        visibleFiles: ['App.tsx'],
      }}
      style={{ height: '100%' }}
    >
      <SandpackLayout style={{ height: '100%' }}>
        <SandpackPreview style={{ height: '100%' }} />
      </SandpackLayout>
    </SandpackProvider>
  );
});
```

**Template Files** (`src/features/Portal/Artifacts/Body/Renderer/React/template.ts`):

```typescript
export const createTemplateFiles = ({ title } = {}) => ({
  'index.html': `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>${title || 'Artifacts App'}</title>
    <script src="https://cdn.tailwindcss.com"></script>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/index.tsx"></script>
  </body>
</html>`,
  
  'vite.config.ts': {
    code: `import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@/components/ui': '@lshay/ui/components/default',
    },
  },
});`,
  },
});
```

### 9. HTML Renderer

**File**: `src/features/Portal/Artifacts/Body/Renderer/HTML.tsx`

Safely renders HTML in an isolated iframe.

```typescript
const HTMLRenderer = memo<{ htmlContent: string; width?: string; height?: string }>(
  ({ htmlContent, width = '100%', height = '100%' }) => {
    const iframeRef = useRef<HTMLIFrameElement>(null);

    useEffect(() => {
      if (!iframeRef.current) return;

      const doc = iframeRef.current.contentDocument;
      if (!doc) return;

      // Write HTML to iframe document
      doc.open();
      doc.write(htmlContent);
      doc.close();
    }, [htmlContent]);

    return (
      <iframe
        ref={iframeRef}
        style={{ border: 'none', width, height }}
        title="html-renderer"
      />
    );
  }
);
```

**Security**: Iframe provides sandboxing to prevent malicious code from affecting the main app.

### 10. SVG Renderer

**File**: `src/features/Portal/Artifacts/Body/Renderer/SVG.tsx`

Renders SVG with interactive features.

```typescript
const SVGRender = memo<{ content: string }>(({ content }) => {
  return (
    <div
      dangerouslySetInnerHTML={{ __html: content }}
      style={{
        width: '100%',
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    />
  );
});
```

## Usage Flow

### Creating an Artifact

1. **User Request**: User asks AI to create substantial content
2. **AI Evaluation**: AI uses `<lobeThinking>` to evaluate if artifact-worthy
3. **AI Response**: AI wraps content in `<lobeArtifact>` tags with attributes
4. **Parsing**: Markdown parser extracts artifact from message
5. **Display**: Inline card appears in chat
6. **Auto-Open**: Portal automatically opens during generation
7. **Rendering**: Content rendered based on type

### Updating an Artifact

1. **User Request**: User asks to modify existing artifact
2. **AI Reuses Identifier**: AI uses same `identifier` attribute
3. **Content Update**: New content replaces old in same artifact
4. **Re-render**: Portal updates to show new content

### User Interaction

1. **Click Card**: User clicks artifact card in chat
2. **Portal Opens**: Side panel opens with artifact viewer
3. **View Modes**: Toggle between code and preview (if applicable)
4. **Close**: Click card again or close button to hide portal

## Display Modes

The artifact viewer supports two display modes:

### Preview Mode (Default)

- Shows rendered output (HTML, React, SVG, etc.)
- Interactive and visual
- Auto-switches when artifact tag closes

### Code Mode

- Shows syntax-highlighted source code
- Useful for copying or understanding implementation
- Always shown for `application/lobe.artifacts.code` type

**Toggle Logic**:

```typescript
const showCode = 
  !isArtifactTagClosed ||              // Still generating
  displayMode === 'code' ||             // User selected code view
  artifactType === 'application/lobe.artifacts.code';  // Non-renderable type
```

## Best Practices

### For AI Assistants

1. **Evaluate First**: Use `<lobeThinking>` to determine if artifact is appropriate
2. **Reuse Identifiers**: Update existing artifacts instead of creating new ones
3. **Complete Content**: Never truncate with "// rest of code remains the same"
4. **Appropriate Types**: Choose correct MIME type for content
5. **Self-Contained**: Ensure artifacts work standalone

### For Developers

1. **Secure Rendering**: Always use iframe for untrusted HTML
2. **Error Handling**: Gracefully handle malformed artifacts
3. **Performance**: Lazy load heavy renderers (React via dynamic import)
4. **Accessibility**: Provide alt text and ARIA labels
5. **Testing**: Test with various content sizes and edge cases

## Security Considerations

### Sandboxing

- **HTML**: Rendered in iframe (isolated DOM)
- **React**: Rendered in Sandpack (isolated environment)
- **SVG**: Sanitized and rendered safely

### Content Restrictions

- **No External Images**: Only placeholder API allowed
- **Limited CDN**: Only `cdnjs.cloudflare.com` for scripts
- **No Arbitrary Code**: React components run in sandboxed environment

### User Protection

- **Visual Indicators**: Clear labeling of artifact source
- **Preview Before Run**: Code view available before execution
- **Error Boundaries**: Catch and display rendering errors

## Performance Optimizations

### 1. Lazy Loading

```typescript
const ReactRenderer = dynamic(() => import('./React'), { ssr: false });
```

Heavy renderers loaded only when needed.

### 2. Memoization

All renderer components use `memo()` to prevent unnecessary re-renders.

### 3. Efficient Parsing

Regex-based parsing is fast and handles streaming content.

### 4. Conditional Rendering

Portal only renders when `showPortal` is true.

## Error Handling

### Malformed Artifacts

- **Missing Attributes**: Use defaults (empty title, 'html' type)
- **Unclosed Tags**: Detect with regex, show partial content
- **Invalid Content**: Display error message in renderer

### Rendering Errors

```typescript
try {
  return <Renderer content={content} type={type} />;
} catch (error) {
  return <ErrorDisplay error={error} />;
}
```

### User Feedback

- **Loading States**: Show spinner during generation
- **Progress Indicators**: Display content length and completion status
- **Error Messages**: Clear, actionable error descriptions

## Testing

### Unit Tests

Test individual components:

```typescript
describe('ArtifactRenderer', () => {
  it('should render HTML artifacts', () => {
    const html = '<h1>Hello</h1>';
    render(<Renderer content={html} type="text/html" />);
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });

  it('should handle malformed content', () => {
    const malformed = '<div>Unclosed';
    render(<Renderer content={malformed} type="text/html" />);
    // Should not crash
  });
});
```

### Integration Tests

Test full artifact flow:

```typescript
describe('Artifact Flow', () => {
  it('should parse and display artifact from message', async () => {
    const message = `
      <lobeArtifact identifier="test" type="text/html" title="Test">
        <h1>Test Content</h1>
      </lobeArtifact>
    `;
    
    render(<ChatMessage content={message} />);
    
    const card = screen.getByText('Test');
    fireEvent.click(card);
    
    await waitFor(() => {
      expect(screen.getByText('Test Content')).toBeInTheDocument();
    });
  });
});
```

## Troubleshooting

### Artifact Not Displaying

**Issue**: Artifact card doesn't appear in chat.

**Solutions**:
1. Check if `<lobeArtifact>` tag is properly formatted
2. Verify all required attributes are present
3. Check browser console for parsing errors
4. Ensure markdown parser is registered

### Portal Not Opening

**Issue**: Clicking artifact card doesn't open portal.

**Solutions**:
1. Check if `openArtifact` action is called
2. Verify portal state in Redux DevTools
3. Check for JavaScript errors in console
4. Ensure portal component is mounted

### Content Not Rendering

**Issue**: Portal opens but content doesn't render.

**Solutions**:
1. Verify artifact type matches content
2. Check renderer component for errors
3. Inspect content extraction from message
4. Test with simpler content

### React Component Errors

**Issue**: React artifacts fail to render.

**Solutions**:
1. Check for syntax errors in component code
2. Verify all imports are from allowed libraries
3. Ensure component has default export
4. Check Sandpack console for errors

## Future Enhancements

### Planned Features

1. **Artifact History**: View previous versions of artifacts
2. **Export Options**: Download artifacts as files
3. **Collaboration**: Share artifacts with other users
4. **Templates**: Pre-built artifact templates
5. **Custom Renderers**: Plugin system for custom types
6. **Offline Support**: Cache artifacts for offline viewing
7. **Search**: Search within artifact content
8. **Annotations**: Add comments to artifacts

### Technical Improvements

1. **Streaming Rendering**: Render content as it streams
2. **Better Error Recovery**: Auto-fix common issues
3. **Performance**: Virtual scrolling for large content
4. **Accessibility**: Enhanced screen reader support
5. **Mobile**: Optimized mobile artifact viewer

## Conclusion

The Artifacts plugin is a powerful feature that extends LobeChat's capabilities beyond simple text conversations. By providing a structured way to create, display, and interact with substantial content, it enables users to leverage AI for creating production-ready code, visualizations, and documents. The architecture is modular, secure, and extensible, making it easy to add new artifact types and renderers as needed.

The combination of intelligent AI guidance (via system prompts), robust parsing, flexible rendering, and intuitive UI makes Artifacts a cornerstone feature for power users and developers working with LobeChat.
