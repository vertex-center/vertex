import { Input, InputProps } from "../Input/Input.tsx";
import { forwardRef, Ref } from "react";

type TextFieldProps = Omit<InputProps, "ref"> & {
    containerRef?: Ref<HTMLDivElement>;
};

export const TextField = forwardRef(
    (props: Readonly<TextFieldProps>, inputRef: Ref<HTMLInputElement>) => {
        const { containerRef, ...others } = props;
        return (
            <Input
                ref={inputRef}
                containerRef={containerRef}
                id="id"
                type="text"
                {...others}
            />
        );
    },
);
