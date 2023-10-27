import { Input, InputProps, InputRef } from "../Input/Input.tsx";
import { Checkbox } from "../Checkbox/Checkbox.tsx";
import {
    Children,
    cloneElement,
    forwardRef,
    HTMLAttributes,
    PropsWithChildren,
    ReactNode,
    useState,
} from "react";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon.tsx";
import "./SelectField.sass";
import cx from "classnames";

export type SelectOptionProps<T> = HTMLAttributes<HTMLDivElement> &
    PropsWithChildren<{
        onClick?: (value: T) => void;
        value: T;
        multiple?: boolean;
        selected?: boolean;
    }>;

export function SelectOption<T>(props: Readonly<SelectOptionProps<T>>) {
    const {
        onClick,
        multiple,
        className,
        value,
        children,
        selected,
        ...others
    } = props;

    return (
        <div
            onClick={() => onClick?.(value)}
            className={cx(
                "select-field-option",
                {
                    "select-field-option-multiple": multiple,
                },
                className,
            )}
            {...others}
        >
            {multiple === true && <Checkbox checked={selected} />}
            {children}
        </div>
    );
}

export type SelectFieldRef = InputRef;

export type SelectFieldProps = Omit<InputProps, "onChange" | "value"> &
    PropsWithChildren<{
        multiple?: boolean;
        value?: ReactNode;
        onChange?: (value: unknown) => void;
    }>;

function _SelectField<T>(
    props: Readonly<SelectFieldProps>,
    ref: SelectFieldRef,
) {
    const {
        children,
        multiple,
        onChange: onChangeProp,
        value,
        ...others
    } = props;

    const [show, setShow] = useState<boolean>(false);

    const onChange = (value: T) => {
        if (!multiple) setShow(false);
        onChangeProp?.(value);
    };

    const toggle = () => setShow(!show);

    return (
        <div
            className={cx("select-field", {
                "select-field-opened": show,
            })}
        >
            <Input
                {...others}
                ref={ref}
                as="div"
                onClick={toggle}
                className="select-field-input"
            >
                {value}
                <MaterialIcon
                    className="select-field-icon"
                    icon="expand_more"
                />
            </Input>
            <div className="select-field-values">
                {Children.map(children, (child) => {
                    if (!child) return;
                    // @ts-ignore
                    return cloneElement(child, {
                        onClick: onChange,
                        multiple: multiple,
                    });
                })}
            </div>
            <div className="select-field-overlay" onClick={toggle} />
        </div>
    );
}

export const SelectField = forwardRef(_SelectField);
