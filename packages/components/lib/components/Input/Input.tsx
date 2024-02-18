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
    as?: ElementType;
};

function _Input(props: Readonly<InputProps>, ref: InputRef) {
    const {
        id,
        as,
        className,
        value: _,
        onChange: __,
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

    return (
        <Component
            ref={ref}
            id={id}
            value={value}
            onChange={onChange}
            className={cx("input", className)}
            children={children}
            {...others}
        />
    );
}

export const Input = forwardRef(_Input);
