import { useEffect, useMemo, useState } from "react";
import debounce from "lodash/debounce";

export const useDebounce = <T,>(value: T, delay = 300): T => {
  const [debouncedValue, setDebouncedValue] = useState(value);
  const updateValue = useMemo(() => debounce(setDebouncedValue, delay), [delay]);

  useEffect(() => {
    updateValue(value);
    return () => updateValue.cancel();
  }, [updateValue, value]);

  return debouncedValue;
};
