import {
    ChangeEvent,
    ElementType,
    forwardRef,
    InputHTMLAttributes,
    Ref,
    useState,
} from "react";
import "./Input.sass";
import cx from "classnames";

export type InputRef = Ref<HTMLInputElement>;

export type InputProps = InputHTMLAttributes<HTMLInputElement> & {
    divRef?: Ref<HTMLDivElement>;
    divProps?: InputHTMLAttributes<HTMLDivElement>;
    label?: string;
    description?: string;
    error?: string;
    as?: ElementType;
};

function _Input(props: Readonly<InputProps>, ref: InputRef) {
    const {
        divRef,
        id,
        as,
        required,
        className,
        value: _,
        onChange: __,
        label,
        description,
        error,
        divProps,
        children,
        ...others
    } = props;

    const controlled = props.value !== undefined;
    const [internalValue, setInternalValue] = useState<string>("");

    const value = controlled ? props.value : internalValue;

    const onChange = (e: ChangeEvent<HTMLInputElement>) => {
        props.onChange?.(e);
        if (!controlled) setInternalValue(e.target.value);
    };

    const Component = as ?? "input";

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
            ref={divRef}
            {...divProps}
            className={cx("input", divProps?.className)}
        >
            {label && (
                <label htmlFor={id} className="input-label">
                    {label} {indicator}
                </label>
            )}
            <Component
                ref={ref}
                value={value}
                onChange={onChange}
                className={cx("input-field", className)}
                children={children}
                {...others}
            />
            {description && !error && (
                <div className="input-description">{description}</div>
            )}
            {error && <div className="input-error">{error}</div>}
        </div>
    );
}

export const Input = forwardRef(_Input);
