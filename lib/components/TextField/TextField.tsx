import { Input, InputProps } from "../Input/Input.tsx";
import { Ref } from "react";

type TextFieldProps = InputProps & {
    ref?: Ref<HTMLInputElement>;
    containerRef?: Ref<HTMLDivElement>;
};

export function TextField(props: Readonly<TextFieldProps>) {
    const { containerRef, ref, ...others } = props;
    return (
        <Input
            inputProps={{ ref }}
            containerRef={containerRef}
            id="id"
            type="text"
            {...others}
        />
    );
}
