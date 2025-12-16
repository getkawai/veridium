## 🎯 Alur Query Assistants List

### 1. **Frontend Store (Zustand + SWR)**

```58:72:src/store/discover/slices/assistant/action.ts
useAssistantList: (params = {}) => {
  const locale = globalHelpers.getCurrentLanguage();
  return useSWR(
    ['assistant-list', locale, ...Object.values(params)].filter(Boolean).join('-'),
    async () =>
      discoverService.getAssistantList({
        ...params,
        page: params.page ? Number(params.page) : 1,
        pageSize: params.pageSize ? Number(params.pageSize) : 21,
      }),
    {
      revalidateOnFocus: false,
    },
  );
},
```

- Menggunakan **SWR** untuk caching dan auto-revalidation
- Key: `['assistant-list', locale, ...params]`
- Default: `page: 1`, `pageSize: 21`

### 2. **Client Service Layer**

```56:67:src/services/discover.ts
getAssistantList = async (params: AssistantQueryParams = {}): Promise<AssistantListResponse> => {
  const locale = globalHelpers.getCurrentLanguage();
  return lambdaClient.market.getAssistantList.query(
    {
      ...params,
      locale,
      page: params.page ? Number(params.page) : 1,
      pageSize: params.pageSize ? Number(params.pageSize) : 20,
    },
    { context: { showNotification: false } },
  );
};
```

### 3. **TRPC Router**

```81:107:src/server/routers/lambda/market/index.ts
getAssistantList: marketProcedure
  .input(
    z
      .object({
        category: z.string().optional(),
        locale: z.string().optional(),
        order: z.enum(['asc', 'desc']).optional(),
        page: z.number().optional(),
        pageSize: z.number().optional(),
        q: z.string().optional(),
        sort: z.nativeEnum(AssistantSorts).optional(),
      })
      .optional(),
  )
  .query(async ({ input, ctx }) => {
    log('getAssistantList input: %O', input);

    try {
      return await ctx.discoverService.getAssistantList(input);
    } catch (error) {
      log('Error fetching assistant list: %O', error);
      throw new TRPCError({
        code: 'INTERNAL_SERVER_ERROR',
        message: 'Failed to fetch assistant list',
      });
    }
  }),
```

### 4. **Server Service (Discover Service)**

```299:416:src/server/services/discover/index.ts
getAssistantList = async (params: AssistantQueryParams = {}): Promise<AssistantListResponse> => {
  log('getAssistantList: params=%O', params);
  const {
    locale,
    category,
    order = 'desc',
    page = 1,
    pageSize = 20,
    q,
    sort = AssistantSorts.CreatedAt,
  } = params;
  let list = await this._getAssistantList(locale);
  const originalCount = list.length;

  if (category) {
    list = list.filter((item) => item.category === category);
    log(
      'getAssistantList: filtered by category "%s", %d -> %d items',
      category,
      originalCount,
      list.length,
    );
  }

  if (q) {
    const beforeFilter = list.length;
    list = list.filter((item) => {
      return [item.author, item.title, item.description, item?.tags]
        .flat()
        .filter(Boolean)
        .join(',')
        .toLowerCase()
        .includes(decodeURIComponent(q).toLowerCase());
    });
    log('getAssistantList: filtered by query "%s", %d -> %d items', q, beforeFilter, list.length);
  }

  if (sort)
```

- Fetch dari **AssistantStore** (external JSON files)
- Filter by: `category`, `q` (search query)
- Sort & pagination

---

## 🎯 Alur Query Sessions List

### 1. **Frontend Store (Zustand + SWR)**

```237:301:src/store/session/slices/session/action.ts
useFetchSessions: (enabled, isLogin) =>
  useClientDataSWR<ChatSessionList>(
    enabled ? [FETCH_SESSIONS_KEY, isLogin] : null,
    () => sessionService.getGroupedSessions(),
    {
      fallbackData: {
        sessionGroups: [],
        sessions: [],
      },
      onSuccess: (data) => {
        if (
          get().isSessionsFirstFetchFinished &&
          isEqual(get().sessions, data.sessions) &&
          isEqual(get().sessionGroups, data.sessionGroups)
        )
          return;

        get().internal_processSessions(
          data.sessions,
          data.sessionGroups,
          n('useFetchSessions/updateData') as any,
        );

        // Sync chat groups from group sessions to chat store
        const groupSessions = data.sessions.filter((session) => session.type === 'group');
        if (groupSessions.length > 0) {
```

- Menggunakan **SWR** dengan `useClientDataSWR`
- Key: `[FETCH_SESSIONS_KEY, isLogin]`
- Memanggil `sessionService.getGroupedSessions()`

### 2. **Client Service Layer**

```55:74:src/services/session/client.ts
getGroupedSessions: ISessionService['getGroupedSessions'] = async () => {
  const { sessions, sessionGroups } = await this.sessionModel.queryWithGroups();
  const chatGroups = await this.chatGroupModel.queryWithMemberDetails();

  const groupSessions = chatGroups.map((group) => {
    const { title, description, avatar, backgroundColor, groupId, ...rest } = group;
    return {
      ...rest,
      group: groupId, // Map groupId to group for consistent API
      meta: { avatar, backgroundColor, description, title },
      type: 'group' as const,
    };
  });

  const allSessions = [...sessions, ...groupSessions].sort(
    (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
  );

  return { sessionGroups, sessions: allSessions };
};
```

### 3. **Database Model (SessionModel)**

```53:78:packages/database/src/models/session.ts
query = async ({ current = 0, pageSize = 9999 } = {}) => {
  const offset = current * pageSize;

  return this.db.query.sessions.findMany({
    limit: pageSize,
    offset,
    orderBy: [desc(sessions.updatedAt)],
    where: and(eq(sessions.userId, this.userId), not(eq(sessions.slug, INBOX_SESSION_ID))),
    with: { agentsToSessions: { columns: {}, with: { agent: true } }, group: true },
  });
};

queryWithGroups = async (): Promise<ChatSessionList> => {
  // 查询所有会话
  const result = await this.query();

  const groups = await this.db.query.sessionGroups.findMany({
    orderBy: [asc(sessionGroups.sort), desc(sessionGroups.createdAt)],
    where: eq(sessions.userId, this.userId),
  });

  return {
    sessionGroups: groups as unknown as ChatSessionList['sessionGroups'],
    sessions: result.map((item) => this.mapSessionItem(item as any)),
  };
};
```

- Query dari **PGLite/PostgreSQL** menggunakan **Drizzle ORM**
- Filter by `userId`
- Exclude `INBOX_SESSION_ID`
- Include relations: `agentsToSessions`, `agent`, `group`
- Order by `updatedAt DESC`

### 4. **TRPC Router (Server-side alternative)**

```99:124:src/server/routers/lambda/session.ts
getGroupedSessions: publicProcedure.query(async ({ ctx }): Promise<ChatSessionList> => {
  if (!ctx.userId) return { sessionGroups: [], sessions: [] };

  const serverDB = await getServerDB();
  const sessionModel = new SessionModel(serverDB, ctx.userId!);
  const chatGroupModel = new ChatGroupModel(serverDB, ctx.userId!);

  const { sessions, sessionGroups } = await sessionModel.queryWithGroups();
  const chatGroups = await chatGroupModel.queryWithMemberDetails();

  const groupSessions: LobeGroupSession[] = chatGroups.map((group) => {
    const { title, description, avatar, backgroundColor, groupId, ...rest } = group;
    return {
      ...rest,
      group: groupId, // Map groupId to group for consistent API
      meta: { avatar, backgroundColor, description, title },
      type: 'group',
    };
  });

  const allSessions = [...sessions, ...groupSessions].sort(
    (a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
  );

  return { sessionGroups, sessions: allSessions };
}),
```

---

## 📊 Perbandingan

| Aspek | Assistants | Sessions |
|-------|-----------|----------|
| **Data Source** | External JSON files (AssistantStore) | Database (PGLite/PostgreSQL) |
| **ORM** | N/A (fetch JSON) | Drizzle ORM |
| **Caching** | SWR | SWR + useClientDataSWR |
| **Filter** | category, search query (q) | userId, exclude INBOX |
| **Pagination** | Yes (page, pageSize) | Yes (current, pageSize) |
| **Relations** | N/A | agentsToSessions, agent, group |
| **Sort** | Multiple sorts (CreatedAt, etc) | updatedAt DESC |

Apakah ada aspek spesifik yang ingin Anda pelajari lebih dalam?