import ChangelogModal from '@/features/ChangelogModal';
// import { ChangelogService } from '@/server/services/changelog';

const Changelog = async () => {
  // const service = new ChangelogService();
  // const id = await service.getLatestChangelogId();

  // Dummy id for UI purposes
  const id = 'dummy-changelog-id';

  return <ChangelogModal currentId={id} />;
};

export default Changelog;
