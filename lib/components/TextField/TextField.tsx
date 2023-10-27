import { Input, InputProps } from "../Input/Input.tsx";
import { forwardRef, Ref } from "react";

type TextFieldRef = Ref<HTMLInputElement>;

type TextFieldProps<T> = InputProps<T> & {
    ref?: Ref<HTMLInputElement>;
    divRef?: Ref<HTMLDivElement>;
};

function _TextField<T>(props: Readonly<TextFieldProps<T>>, ref: TextFieldRef) {
    const { divRef, ...others } = props;
    return (
        <Input
            divRef={divRef}
            inputProps={{ ref }}
            id="id"
            type="text"
            {...others}
        />
    );
}

export const TextField = forwardRef(_TextField);
