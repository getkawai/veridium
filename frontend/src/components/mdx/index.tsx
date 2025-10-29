import { Typography as Typo, TypographyProps } from '@lobehub/ui';

export const Typography = ({
  children,
  mobile,
  style,
  ...rest
}: { mobile?: boolean } & TypographyProps) => {
  const headerMultiple = mobile ? 0.2 : 0.4;
  return (
    <Typo
      fontSize={14}
      headerMultiple={headerMultiple}
      style={{ width: '100%', ...style }}
      {...rest}
    >
      {children}
    </Typo>
  );
};