import { Input, InputProps, InputRef } from "../Input/Input.tsx";
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
import { Check } from "@phosphor-icons/react";

export type SelectOptionProps<T> = HTMLAttributes<HTMLDivElement> &
    PropsWithChildren<{
        onClick?: (value: T) => void;
        value: T;
        left?: ReactNode;
    }>;

export function SelectOption<T>(props: Readonly<SelectOptionProps<T>>) {
    const { onClick, className, value, children, left, ...others } = props;

    return (
        <div
            onClick={() => onClick?.(value)}
            className={cx("select-field-option", className)}
            {...others}
        >
            <div className="select-field-option-left">{left}</div>
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
        value?: T;
        valueRender?: (value?: T) => ReactNode;
        leftIcon?: ReactNode;
        onChange?: (value: T) => void;
        filter?: (value: T, search: string) => boolean;
        textNoResults?: string;
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
        valueRender,
        leftIcon,
        filter,
        textNoResults = "No results found.",
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

    const items = Children.map(children, (child) => {
        if (!child) return null;
        // @ts-expect-error props are too hard to type
        const v = child.props.value;
        if (filter && !filter(v, search)) return null;
        // @ts-expect-error cloneElement is too hard to type
        return cloneElement(child, {
            onClick: onChange,
            multiple: multiple,
            left: v === value && <Check />,
        });
    })?.filter((v) => v);

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
                {leftIcon}
                {valueRender?.(value) || value?.toString()}
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
                    {items}
                    {(!items || items?.length === 0) && (
                        <div className="select-field-no-results">
                            {textNoResults}
                        </div>
                    )}
                </div>
            </div>
            <div className="select-field-overlay" onClick={toggle} />
        </div>
    );
}

export const SelectField = forwardRef(_SelectField);
