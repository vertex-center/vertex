import { Input, InputProps } from "../Input/Input.tsx";
import cx from "classnames";
import { HTMLProps } from "react";
import "./TextField.sass";

type TextFieldProps = HTMLProps<HTMLDivElement> & {
    inputProps?: InputProps;

    label?: string;
};

export function TextField(props: Readonly<TextFieldProps>) {
    const { inputProps, className, label, ...others } = props;

    return (
        <div className={cx("text-field", className)} {...others}>
            <label className="text-field-label">{label}</label>
            <Input {...inputProps} />
        </div>
    );
}
