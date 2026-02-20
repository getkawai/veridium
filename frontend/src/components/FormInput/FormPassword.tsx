import { InputPassword, InputProps as Props } from '@lobehub/ui';
import { InputRef } from 'antd/es/input/Input';
import { memo, useRef, useState } from 'react';

interface FormPasswordProps extends Omit<Props, 'onChange'> {
  onChange?: (value: string) => void;
}

const FormPassword = memo<FormPasswordProps>(({ onChange, value: defaultValue, ...props }) => {
  const ref = useRef<InputRef>(null);
  const isChineseInput = useRef(false);

  const [value, setValue] = useState(defaultValue as string);
  const prevDefaultValue = useRef(defaultValue);
  if (prevDefaultValue.current !== defaultValue) {
    prevDefaultValue.current = defaultValue;
    setValue(defaultValue as string);
  }

  return (
    <InputPassword
      onBlur={() => {
        onChange?.(value);
      }}
      onChange={(e) => {
        setValue(e.target.value);
      }}
      onCompositionEnd={() => {
        isChineseInput.current = false;
      }}
      onCompositionStart={() => {
        isChineseInput.current = true;
      }}
      onPressEnter={() => {
        if (isChineseInput.current) return;
        onChange?.(value);
      }}
      ref={ref}
      {...props}
      value={value}
    />
  );
});

FormPassword.displayName = 'FormPassword';

export default FormPassword;
