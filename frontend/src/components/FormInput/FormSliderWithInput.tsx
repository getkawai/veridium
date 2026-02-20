import { SliderWithInput, SliderWithInputProps } from '@lobehub/ui';
import { memo, useRef, useState } from 'react';

interface FormSliderWithInputProps extends Omit<SliderWithInputProps, 'onChange' | 'value'> {
  onChange?: (value: number) => void;
  value?: number;
}

/**
 * Form-integrated slider with delayed onChange behavior.
 * Only triggers onChange on blur to prevent excessive updates during user interaction.
 */
const FormSliderWithInput = memo<FormSliderWithInputProps>(
  ({ onChange, value: defaultValue, ...props }) => {
    const [value, setValue] = useState(defaultValue ?? 0);
    const prevDefaultValue = useRef(defaultValue);
    if (prevDefaultValue.current !== defaultValue) {
      prevDefaultValue.current = defaultValue;
      setValue(defaultValue ?? 0);
    }

    return (
      <SliderWithInput
        onBlur={() => {
          onChange?.(value);
        }}
        onChange={(newValue) => {
          if (typeof newValue === 'number') {
            setValue(newValue);
          }
        }}
        {...props}
        value={value}
      />
    );
  },
);

FormSliderWithInput.displayName = 'FormSliderWithInput';

export default FormSliderWithInput;
