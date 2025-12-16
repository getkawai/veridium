## How Message Translation Works in LobeChat

The message translation feature in LobeChat uses AI to translate messages between different languages. Here's how it works:

### 1. **Translation Prompt Chain** (`chainTranslate`)

The translation is powered by a prompt template defined in:

```11:24:packages/prompts/src/chains/translate.ts
Rules:
- Output ONLY the translated text, no explanations or additional context
- Preserve technical terms, code identifiers, API keys, and proper nouns exactly as they appear
- Maintain the original formatting and structure
- Use natural, idiomatic expressions in the target language`,
      role: 'system',
    },
    {
      content,
      role: 'user',
    },
  ],
});
```

This creates a chat payload with:
- A **system message** instructing the AI to act as a professional translator
- The **user's message content** to be translated
- The **target language** specified

### 2. **Translation Action Flow**

The main translation logic is in the Zustand store action:

```43:96:src/store/chat/slices/translate/action.ts
  translateMessage: async (id, targetLang) => {
    const { internal_toggleChatLoading, updateMessageTranslate, internal_dispatchMessage } = get();

    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    // Get current agent for translation
    const translationSetting = systemAgentSelectors.translation(useUserStore.getState());

    // create translate extra
    await updateMessageTranslate(id, { content: '', from: '', to: targetLang });

    internal_toggleChatLoading(true, id, n('translateMessage(start)', { id }));

    let content = '';
    let from = '';

    // detect from language
    chatService.fetchPresetTaskResult({
      onFinish: async (data) => {
        if (data && supportLocales.includes(data)) from = data;

        await updateMessageTranslate(id, { content, from, to: targetLang });
      },
      params: merge(translationSetting, chainLangDetect(message.content)),
      trace: get().getCurrentTracePayload({ traceName: TraceNameMap.LanguageDetect }),
    });

    // translate to target language
    await chatService.fetchPresetTaskResult({
      onFinish: async (content) => {
        await updateMessageTranslate(id, { content, from, to: targetLang });
        internal_toggleChatLoading(false, id);
      },
      onMessageHandle: (chunk) => {
        switch (chunk.type) {
          case 'text': {
            internal_dispatchMessage({
              id,
              key: 'translate',
              type: 'updateMessageExtra',
              value: produce({ content: '', from, to: targetLang }, (draft) => {
                content += chunk.text;
                draft.content += content;
              }),
            });
            break;
          }
        }
      },
      params: merge(translationSetting, chainTranslate(message.content, targetLang)),
      trace: get().getCurrentTracePayload({ traceName: TraceNameMap.Translator }),
    });
  },
```

The process involves **two AI calls**:

#### Step 1: Language Detection
- Uses `chainLangDetect` to identify the source language
- The AI returns a locale code (e.g., `en-US`, `zh-CN`)
- This is stored in the `from` field

#### Step 2: Translation
- Uses `chainTranslate` with the message content and target language
- Streams the translated text in real-time using `onMessageHandle`
- Updates the UI progressively as chunks arrive
- Stores the final translation in the database

### 3. **Database Storage**

Translations are stored in a separate table:

```189:209:packages/database/src/schemas/message.ts
export const messageTranslates = pgTable(
  'message_translates',
  {
    id: text('id')
      .references(() => messages.id, { onDelete: 'cascade' })
      .primaryKey(),
    content: text('content'),
    from: text('from'),
    to: text('to'),
    clientId: text('client_id'),
    userId: text('user_id')
      .references(() => users.id, { onDelete: 'cascade' })
      .notNull(),
  },
  (t) => ({
    clientIdUnique: uniqueIndex('message_translates_client_id_user_id_unique').on(
      t.clientId,
      t.userId,
    ),
  }),
);
```

This table stores:
- `id`: References the original message
- `content`: The translated text
- `from`: Source language (detected)
- `to`: Target language (user-selected)
- `userId`: Owner of the translation

### 4. **UI Display**

The translation is displayed using the `Translate` component:

```18:69:src/features/Conversation/components/Extras/Translate.tsx
const Translate = memo<TranslateProps>(({ content = '', from, to, id, loading }) => {
  const theme = useTheme();
  const { t } = useTranslation('common');
  const [show, setShow] = useState(true);
  const clearTranslate = useChatStore((s) => s.clearTranslate);

  const { message } = App.useApp();
  return (
    <Flexbox gap={8}>
      <Flexbox align={'center'} horizontal justify={'space-between'}>
        <div>
          <Flexbox gap={4} horizontal>
            <Tag style={{ margin: 0 }}>{from ? t(`lang.${from}` as any) : '...'}</Tag>
            <Icon color={theme.colorTextTertiary} icon={ChevronsRight} />
            <Tag>{t(`lang.${to}` as any)}</Tag>
          </Flexbox>
        </div>
        <Flexbox horizontal>
          <ActionIcon
            icon={CopyIcon}
            onClick={async () => {
              await copyToClipboard(content);
              message.success(t('copySuccess'));
            }}
            size={'small'}
            title={t('copy')}
          />
          <ActionIcon
            icon={TrashIcon}
            onClick={() => {
              clearTranslate(id);
            }}
            size={'small'}
            title={t('translate.clear', { ns: 'chat' })}
          />
          <ActionIcon
            icon={show ? ChevronDown : ChevronUp}
            onClick={() => {
              setShow(!show);
            }}
            size={'small'}
          />
        </Flexbox>
      </Flexbox>
      {!show ? null : loading && !content ? (
        <BubblesLoading />
      ) : (
        <Markdown variant={'chat'}>{content}</Markdown>
      )}
    </Flexbox>
  );
});
```

This component shows:
- Source and target language tags (e.g., `en-US` → `zh-CN`)
- Copy, delete, and collapse/expand buttons
- Loading indicator while translating
- The translated content in Markdown format

### 5. **User Interaction**

Users can trigger translation from the message actions menu in either assistant or user messages. When clicked, it calls `translateMessage(id, lang)` with the selected target language.

### Summary

The translation workflow:
1. User clicks translate and selects target language
2. System detects source language using AI
3. System translates content using AI (streaming)
4. Translation is stored in database
5. UI displays translation with language tags and controls

The system uses AI models configured in the user's translation settings, making it flexible and supporting any language the underlying model supports.