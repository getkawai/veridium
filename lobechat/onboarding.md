# Onboarding and Welcome Experience Documentation

This document describes the onboarding and welcome experience in LobeChat, covering the components, state management, and localization.

## 1. Overview

The onboarding flow ensures new users are properly introduced to the application. It consists of a redirection logic for first-time users and a welcome interface displayed in the chat workspace.

## 2. Redirection Logic

When a user is not yet "onboarded," the system redirects them to an onboarding path (typically `/onboard`).

### Key Files:
- [StoreInitialization.tsx](file:///Users/yuda/github.com/lobehub/lobe-chat/src/layout/GlobalProvider/StoreInitialization.tsx): Checks the `isOnboard` state during initialization.
  ```typescript
  useInitUserState(isLoginOnInit, serverConfig, {
    onSuccess: (state) => {
      if (state.isOnboard === false) {
        router.push('/onboard');
      }
    },
  });
  ```
- [Redirect.tsx](file:///Users/yuda/github.com/lobehub/lobe-chat/src/app/[variants]/loading/Server/Redirect.tsx): Handles redirection during server-side loading.
  ```typescript
  if (!isOnboard) {
    router.replace('/onboard');
    return;
  }
  ```

> [!NOTE]
> In the community version, `isOnboard` is often defaulted to `true` in the backend to reduce friction.

## 3. Onboarding State Management

The onboarding status is managed within the `user` store.

### Key Files:
- [action.ts](file:///Users/yuda/github.com/lobehub/lobe-chat/src/store/user/slices/common/action.ts): Defines the `isOnboard` state in the `common` slice.
- [initialState.ts](file:///Users/yuda/github.com/lobehub/lobe-chat/src/store/user/slices/common/initialState.ts): Sets the default value (usually `false`).
- [user.ts (Server Router)](file:///Users/yuda/github.com/lobehub/lobe-chat/src/server/routers/lambda/user.ts): Provides the `makeUserOnboarded` mutation to update the status in the database.

## 4. Welcome Components

The welcome experience is integrated into the chat interface via the `WelcomeChatItem`.

### Component Hierarchy:
- [WelcomeChatItem](file:///Users/yuda/github.com/lobehub/lobe-chat/src/app/[variants]/(main)/chat/(workspace)/@conversation/features/ChatList/WelcomeChatItem/index.tsx): A wrapper that decides whether to show `AgentWelcome` or `GroupWelcome`.
- [AgentWelcome](file:///Users/yuda/github.com/lobehub/lobe-chat/src/app/[variants]/(main)/chat/(workspace)/@conversation/features/ChatList/WelcomeChatItem/AgentWelcome/index.tsx): The default greeting for individual AI agents.
- [GroupWelcome](file:///Users/yuda/github.com/lobehub/lobe-chat/src/app/[variants]/(main)/chat/(workspace)/@conversation/features/ChatList/WelcomeChatItem/GroupWelcome/index.tsx): The greeting for group chat sessions.
- [OpeningQuestions](file:///Users/yuda/github.com/lobehub/lobe-chat/src/app/[variants]/(main)/chat/(workspace)/@conversation/features/ChatList/WelcomeChatItem/AgentWelcome/OpeningQuestions.tsx): Displays interactive suggested questions to help users start the conversation.

## 5. Localization

Onboarding strings are localized using i18next.

- [welcome.ts](file:///Users/yuda/github.com/lobehub/lobe-chat/src/locales/default/welcome.ts): Contains all greeting messages, slogans, and guide content.
  - `guide.defaultMessage`: The main introductory text.
  - `slogan`: Brand slogans like "给自己一个更聪明的大脑".
  - `welcome`: Time-based greetings (Morning, Afternoon, etc.).
