import { Input, InputProps } from "../Input/Input.tsx";
import { Ref } from "react";

type TextFieldProps<T> = InputProps<T> & {
    ref?: Ref<HTMLInputElement>;
    containerRef?: Ref<HTMLDivElement>;
};

export function TextField<T>(props: Readonly<TextFieldProps<T>>) {
    const { containerRef, ref, ...others } = props;
    return (
        <Input
            ref={containerRef}
            inputProps={{ ref }}
            id="id"
            type="text"
            {...others}
        />
    );
}
