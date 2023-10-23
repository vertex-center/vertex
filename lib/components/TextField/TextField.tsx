import { Input, InputProps } from "../Input/Input.tsx";
import cx from "classnames";
import { HTMLProps } from "react";
import "./TextField.sass";

type TextFieldProps = HTMLProps<HTMLDivElement> & {
    inputProps?: InputProps;

    required?: boolean;
    label?: string;
    description?: string;
    error?: string;
};

export function TextField(props: Readonly<TextFieldProps>) {
    const {
        inputProps,
        className,
        required,
        label,
        description,
        error,
        ...others
    } = props;

    let indicator;
    if (required) {
        indicator = <span className="text-field-required">*</span>;
    } else {
        indicator = <span className="text-field-optional">(optional)</span>;
    }

    return (
        <div className={cx("text-field", className)} {...others}>
            {label && (
                <label className="text-field-label">
                    {label} {indicator}
                </label>
            )}
            <Input type="text" {...inputProps} />
            {description && !error && (
                <div className="text-field-description">{description}</div>
            )}
            {error && <div className="text-field-error">{error}</div>}
        </div>
    );
}
