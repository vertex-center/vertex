import {
    ElementType,
    HTMLAttributes,
    HTMLInputTypeAttribute,
    HTMLProps,
    Ref,
} from "react";
import "./Input.sass";
import cx from "classnames";

export type InputProps = HTMLAttributes<HTMLDivElement> & {
    label?: string;
    description?: string;
    error?: string;
    as?: ElementType;
    containerRef?: Ref<HTMLDivElement>;
    inputProps?: HTMLProps<HTMLInputElement>;
    required?: boolean;
    disabled?: boolean;
    type?: HTMLInputTypeAttribute;
};

export function Input(props: Readonly<InputProps>) {
    const {
        className,
        id,
        as,
        required,
        placeholder,
        disabled,
        type,
        label,
        description,
        error,
        inputProps,
        children,
        ...others
    } = props;

    const Component = as || "input";

    if (!id) {
        console.warn("Input is missing an id", { label, description });
    }

    let indicator;
    if (required) {
        indicator = <span className="input-required">*</span>;
    } else {
        indicator = <span className="input-optional">(optional)</span>;
    }

    return (
        <div className={cx("input", className)} {...others}>
            {label && (
                <label htmlFor={id} className="input-label">
                    {label} {indicator}
                </label>
            )}
            <Component
                id={id}
                placeholder={placeholder}
                required={required}
                disabled={disabled}
                type={type}
                {...inputProps}
                className={cx("input-field", inputProps?.className)}
                children={children}
            />
            {description && !error && (
                <div className="input-description">{description}</div>
            )}
            {error && <div className="input-error">{error}</div>}
        </div>
    );
}
