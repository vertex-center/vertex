import { forwardRef, HTMLProps, Ref } from "react";
import "./Input.sass";
import cx from "classnames";

export type InputProps = HTMLProps<HTMLDivElement> & {
    label?: string;
    description?: string;
    error?: string;
    containerRef?: Ref<HTMLDivElement>;
    inputProps?: HTMLProps<HTMLInputElement>;
};

export const Input = forwardRef(
    (props: Readonly<InputProps>, inputRef: Ref<HTMLInputElement>) => {
        const {
            className,
            id,
            ref: containerRef,
            required,
            placeholder,
            disabled,
            label,
            description,
            error,
            inputProps,
            ...others
        } = props;

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
            <div
                ref={containerRef}
                className={cx("input", className)}
                {...others}
            >
                {label && (
                    <label htmlFor={id} className="input-label">
                        {label} {indicator}
                    </label>
                )}
                <input
                    ref={inputRef}
                    id={id}
                    placeholder={placeholder}
                    required={required}
                    disabled={disabled}
                    {...inputProps}
                    className={cx("input-field", inputProps?.className)}
                />
                {description && !error && (
                    <div className="input-description">{description}</div>
                )}
                {error && <div className="input-error">{error}</div>}
            </div>
        );
    },
);
