import React, { memo, useMemo, useState } from 'react';
import { shallow } from 'zustand/shallow';

import { DEFAULT_AVATAR } from '@/const/meta';
import { isDesktop } from '@/const/version';
import { useChatStore } from '@/store/chat';
import { chatSelectors } from '@/store/chat/selectors';
import { useGlobalStore } from '@/store/global';
import { useSessionStore } from '@/store/session';
import { sessionSelectors } from '@/store/session/selectors';
import { useUserStore } from '@/store/user';
import { userProfileSelectors } from '@/store/user/selectors';
import { getNullableString, Session } from '@/types/database';
import { LobeSessionType } from '@/types/session';

import ListItem from '../../ListItem';
import CreateGroupModal from '../../Modals/CreateGroupModal';
import Actions from './Actions';

interface SessionItemProps {
  id: string;
}

const SessionItem = memo<SessionItemProps>(({ id }) => {
  const [open, setOpen] = useState(false);
  const [createGroupModalOpen, setCreateGroupModalOpen] = useState(false);

  const openSessionInNewWindow = useGlobalStore((s) => s.openSessionInNewWindow);

  const [active] = useSessionStore((s) => [s.activeId === id]);
  const [loading] = useChatStore((s) => [chatSelectors.isAIGenerating(s) && id === s.activeId]);

  const [pin, title, description, avatar, avatarBackground, updateAt, members, sessionGroup, sessionType] =
    useSessionStore((s) => {
      const session: Session = sessionSelectors.getSessionById(id)(s);
      if (!session) return [false, '', '', DEFAULT_AVATAR, undefined, undefined, [] as Array<{ avatar: string; backgroundColor?: string }>, undefined, 'agent'];

      // Get metadata from session directly (no nested meta object)
      const sessionTitle = getNullableString(session.title);
      const sessionDescription = getNullableString(session.description);
      const sessionAvatar = getNullableString(session.avatar) || DEFAULT_AVATAR;
      const sessionBg = getNullableString(session.backgroundColor);
      const sessionGroupId = getNullableString(session.groupId);
      const sessionTypeStr = getNullableString(session.type) as LobeSessionType || 'agent';

      return [
        Boolean(session.pinned),
        sessionTitle || '',
        sessionDescription || '',
        sessionAvatar,
        sessionBg,
        session.updatedAt,
        [] as Array<{ avatar: string; backgroundColor?: string }>, // members - would need to be fetched separately for group sessions
        sessionGroupId,
        sessionTypeStr,
      ];
    });

  const handleDoubleClick = () => {
    if (isDesktop) {
      openSessionInNewWindow(id);
    }
  };

  const handleDragStart = (e: React.DragEvent) => {
    // Set drag data to identify the session being dragged
    e.dataTransfer.setData('text/plain', id);
  };

  const handleDragEnd = (e: React.DragEvent) => {
    // If drag ends without being dropped in a valid target, open in new window
    if (isDesktop && e.dataTransfer.dropEffect === 'none') {
      openSessionInNewWindow(id);
    }
  };

  const actions = useMemo(
    () => (
      <Actions
        group={sessionGroup}
        id={id}
        openCreateGroupModal={() => setCreateGroupModalOpen(true)}
        parentType={sessionType as 'agent' | 'group'}
        setOpen={setOpen}
      />
    ),
    [sessionGroup, id],
  );

  const addon = useMemo(
    () =>
      description ? (
        <div
          style={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}
        >
          {description}
        </div>
      ) : undefined,
    [description],
  );

  const currentUser = useUserStore((s) => ({
    avatar: userProfileSelectors.userAvatar(s),
    name: userProfileSelectors.displayUserName(s) || userProfileSelectors.nickName(s) || 'You',
  }));

  const sessionAvatar: string | { avatar: string; background?: string }[] =
    sessionType === 'group'
      ? [
        {
          avatar: currentUser.avatar || DEFAULT_AVATAR,
          background: undefined,
        },
        ...(members?.map((member) => ({
          avatar: member.avatar || DEFAULT_AVATAR,
          background: member.backgroundColor || undefined,
        })) || []),
      ]
      : avatar;

  return (
    <>
      <ListItem
        actions={actions}
        active={active}
        addon={addon}
        avatar={sessionAvatar as any} // Fix: Bypass complex intersection type ReactNode & avatar type
        avatarBackground={avatarBackground}
        date={updateAt?.valueOf()}
        draggable={isDesktop}
        key={id}
        loading={loading}
        onDoubleClick={handleDoubleClick}
        onDragEnd={handleDragEnd}
        onDragStart={handleDragStart}
        pin={pin}
        showAction={open}
        styles={{
          container: {
            gap: 12,
          },
          content: {
            gap: 6,
            maskImage: `linear-gradient(90deg, #000 90%, transparent)`,
          },
        }}
        title={title}
        type={sessionType as 'agent' | 'group' | 'inbox' | undefined}
      />
      <CreateGroupModal
        id={id}
        onCancel={() => setCreateGroupModalOpen(false)}
        open={createGroupModalOpen}
      />
    </>
  );
}, shallow);

export default SessionItem;
