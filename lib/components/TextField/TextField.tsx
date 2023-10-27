import { Input, InputProps } from "../Input/Input.tsx";
import { forwardRef, Ref } from "react";

type TextFieldRef = Ref<HTMLInputElement>;

type TextFieldProps = InputProps & {
    ref?: Ref<HTMLInputElement>;
    divRef?: Ref<HTMLDivElement>;
};

function _TextField(props: Readonly<TextFieldProps>, ref: TextFieldRef) {
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
