import {
    ChangeEvent,
    ElementType,
    forwardRef,
    HTMLAttributes,
    HTMLInputTypeAttribute,
    HTMLProps,
    Ref,
    useEffect,
    useState,
} from "react";
import "./Input.sass";
import cx from "classnames";

export type InputRef = Ref<HTMLInputElement>;

export type InputProps<T> = Omit<HTMLAttributes<HTMLDivElement>, "onChange"> & {
    divRef?: Ref<HTMLDivElement>;
    value?: T;
    onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
    label?: string;
    description?: string;
    error?: string;
    as?: ElementType;
    inputProps?: HTMLProps<HTMLInputElement>;
    required?: boolean;
    disabled?: boolean;
    type?: HTMLInputTypeAttribute;
};

function _Input<T>(props: Readonly<InputProps<T>>, ref: InputRef) {
    const {
        divRef,
        className,
        id,
        as,
        required,
        placeholder,
        disabled,
        value: initialValue,
        onChange: _,
        type,
        label,
        description,
        error,
        inputProps,
        children,
        ...others
    } = props;

    const [value, setValue] = useState<T | undefined>(initialValue);

    useEffect(() => {
        setValue(initialValue);
    }, [initialValue]);

    const onChange = (e: ChangeEvent<HTMLInputElement>) => {
        setValue(e.target.value as T);
        props.onChange?.(e);
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
        <div ref={divRef} className={cx("input", className)} {...others}>
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
                value={value}
                onChange={onChange}
                type={type}
                {...inputProps}
                ref={ref}
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

export const Input = forwardRef(_Input);
