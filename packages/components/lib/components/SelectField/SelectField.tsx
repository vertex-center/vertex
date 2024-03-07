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
    ChangeEvent,
    Fragment,
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

export function SelectSearch(props: InputProps) {
    return (
        <Input
            className="select-field-search"
            placeholder="Search..."
            {...props}
        />
    );
}

export function SelectDivider() {
    return <div className="select-field-divider" />;
}

export type SelectFieldRef = InputRef;

export type SelectFieldProps<T> = Omit<InputProps, "onChange" | "value"> &
    PropsWithChildren<{
        multiple?: boolean;
        value?: ReactNode;
        onChange?: (value: T) => void;
        filter?: (value: T, search: string) => boolean;
    }>;

function _SelectField<T>(
    props: Readonly<SelectFieldProps<T>>,
    ref: SelectFieldRef,
) {
    const {
        children,
        multiple,
        onChange: onChangeProp,
        value,
        filter,
        ...others
    } = props;

    const [show, setShow] = useState<boolean>(false);
    const [search, setSearch] = useState<string>("");

    const onChange = (value: T) => {
        if (!multiple) setShow(false);
        onChangeProp?.(value);
    };

    const onSearchChange = (e: ChangeEvent<HTMLInputElement>) => {
        setSearch(e.currentTarget.value);
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
                onClick={() => setShow(true)}
                className="select-field-input"
            >
                {value}
                <MaterialIcon
                    className="select-field-icon"
                    icon="expand_more"
                />
            </Input>
            <div className="select-field-dropdown">
                {filter && (
                    <Fragment>
                        <SelectSearch
                            value={search}
                            onChange={onSearchChange}
                        />
                        <SelectDivider />
                    </Fragment>
                )}
                <div className="select-field-values">
                    {Children.map(children, (child) => {
                        if (!child) return;
                        // @ts-expect-error props are too hard to type
                        if (filter && !filter(child.props.value, search))
                            return;
                        // @ts-expect-error cloneElement is too hard to type
                        return cloneElement(child, {
                            onClick: onChange,
                            multiple: multiple,
                        });
                    })}
                </div>
            </div>
            <div className="select-field-overlay" onClick={toggle} />
        </div>
    );
}

export const SelectField = forwardRef(_SelectField);
