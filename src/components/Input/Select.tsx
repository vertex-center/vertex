import React, {
    Children,
    cloneElement,
    Fragment,
    HTMLProps,
    useState,
} from "react";
import classNames from "classnames";

import styles from "./Input.module.sass";
import Spacer from "../Spacer/Spacer";
import Icon from "../Icon/Icon";
import Checkbox from "../Checkbox/Checkbox";

type OptionProps = HTMLProps<HTMLDivElement> & {
    onClick?: (value: any) => void;
    value?: any;
};

export function SelectOption(props: Readonly<OptionProps>) {
    const {
        className,
        selected,
        multiple,
        onClick,
        value,
        children,
        ...others
    } = props;

    return (
        <Fragment>
            <div
                onClick={() => onClick(value)}
                className={classNames({
                    [styles.selectItem]: true,
                    [styles.selectItemMultiple]: multiple,
                    [className]: true,
                })}
                {...others}
            >
                {multiple === true && <Checkbox checked={selected} />}
                {children}
            </div>
        </Fragment>
    );
}

type SelectValueProps = HTMLProps<HTMLDivElement>;

export function SelectValue(props: Readonly<SelectValueProps>) {
    const { className, ...others } = props;

    return (
        <div
            className={classNames(styles.selectValue, className)}
            {...others}
        />
    );
}

type SelectOptionsProps = HTMLProps<HTMLDivElement> & {
    opened?: boolean;
};

export function SelectOptions(props: Readonly<SelectOptionsProps>) {
    const { opened, className, ...others } = props;

    return (
        <div
            className={classNames({
                [styles.selectItems]: true,
                [styles.selectItemsOpened]: opened,
                [className]: true,
            })}
            {...others}
        />
    );
}

type Props = HTMLProps<HTMLDivElement> & {
    value?: Element;
    multiple?: boolean;
    onChange?: (value: any) => void;
};

export default function Select(props: Readonly<Props>) {
    const { className, value, label, children, disabled, required, ...others } =
        props;

    const [opened, setOpened] = useState(false);

    const toggle = () => {
        setOpened((o) => {
            if (!o && props.disabled) {
                // Don't open if the component is disabled
                return o;
            }
            return !o;
        });
    };

    const onChange = (value: any) => {
        if (!props.multiple) {
            toggle();
        }
        props.onChange?.(value);
    };

    return (
        <div className={styles.select}>
            {label && (
                <label className={styles.label}>
                    {label}
                    {required && <span className={styles.required}>*</span>}
                    {!required && (
                        <span className={styles.optional}>(optional)</span>
                    )}
                </label>
            )}
            <div
                className={classNames({
                    [styles.input]: true,
                    [styles.inputDisabled]: disabled,
                    [className]: true,
                })}
                {...others}
                onClick={toggle}
            >
                {value}
                <Spacer />
                {opened && <Icon name="expand_less" />}
                {!opened && <Icon name="expand_more" />}
            </div>
            <SelectOptions opened={opened}>
                {Children.map(children, (child) => {
                    if (!child) return;
                    // @ts-ignore
                    return cloneElement(child, {
                        onClick: onChange,
                    });
                })}
            </SelectOptions>
            <div
                className={classNames({
                    [styles.selectOverlay]: true,
                    [styles.selectOverlayOpened]: opened,
                })}
                onClick={toggle}
            />
        </div>
    );
}
