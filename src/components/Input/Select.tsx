import React, { Children, cloneElement, HTMLProps, useState } from "react";
import classNames from "classnames";

import styles from "./Input.module.sass";
import Spacer from "../Spacer/Spacer";
import Symbol from "../Symbol/Symbol";

type OptionProps = HTMLProps<HTMLDivElement> & {
    onClick?: (value: any) => void;
    value?: any;
};

export function SelectOption(props: OptionProps) {
    const { className, onClick, value, ...others } = props;

    return (
        <div
            onClick={() => onClick(value)}
            className={classNames(styles.selectItem, className)}
            {...others}
        />
    );
}

type SelectValueProps = HTMLProps<HTMLDivElement>;

export function SelectValue(props: SelectValueProps) {
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
    toggle?: () => void;
};

export function SelectOptions(props: SelectOptionsProps) {
    const { opened, toggle, className, ...others } = props;

    return (
        <div
            onClick={toggle}
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
    onChange?: (value: any) => void;
};

export default function Select(props: Props) {
    const { className, value, label, children, ...others } = props;

    const [opened, setOpened] = useState(false);

    const toggle = () => setOpened(!opened);

    const onChange = (value: any) => {
        toggle();
        props.onChange?.(value);
    };

    return (
        <div className={styles.select}>
            {label && <label className={styles.label}>{label}</label>}
            <div
                className={classNames(styles.input, className)}
                {...others}
                onClick={toggle}
            >
                {value}
                <Spacer />
                {opened && <Symbol name="expand_less" />}
                {!opened && <Symbol name="expand_more" />}
            </div>
            <SelectOptions opened={opened} toggle={toggle}>
                {Children.map(children, (child) => {
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
