import { Input, InputProps } from "../Input/Input.tsx";
import cx from "classnames";
import { HTMLProps } from "react";
import "./TextField.sass";

type TextFieldProps = HTMLProps<HTMLDivElement> & {
    inputProps?: InputProps;

    label?: string;
    required?: boolean;
};

export function TextField(props: Readonly<TextFieldProps>) {
    const { inputProps, className, label, required, ...others } = props;

    let indicator;
    if (required) {
        indicator = <span className="text-field-required">*</span>;
    } else {
        indicator = <span className="text-field-optional">(optional)</span>;
    }

    return (
        <div className={cx("text-field", className)} {...others}>
            <label className="text-field-label">
                {label} {indicator}
            </label>
            <Input type="text" {...inputProps} />
        </div>
    );
}
