import { Fragment, HTMLProps } from "react";

import styles from "./Button.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";
import Spacer from "../Spacer/Spacer";
import { Link } from "react-router-dom";

export type ButtonProps = HTMLProps<HTMLButtonElement> & {
    leftIcon?: string | JSX.Element;
    rightIcon?: string | JSX.Element;

    selected?: boolean;
    selectable?: boolean;

    loading?: boolean;
    disabled?: boolean;

    onClick?: () => void;
    to?: string;

    color?: "default" | "red";

    // types
    primary?: boolean;
    large?: boolean;
};

export default function Button(props: Readonly<ButtonProps>) {
    const {
        children,
        leftIcon,
        rightIcon,
        loading,
        disabled,
        primary,
        large,
        selected,
        selectable,
        onClick,
        to,
        className,
        type,
        color,
        ...others
    } = props;

    const content = (
        <Fragment>
            {leftIcon && <Icon className={styles.icon} name={leftIcon} />}
            <div>{children}</div>
            {rightIcon && <Icon className={styles.icon} name={rightIcon} />}
            {selectable && (
                <Fragment>
                    <Spacer />
                    <Icon
                        style={{
                            opacity: selected ? 1 : 0.5,
                            color: selected ? "var(--bg-accent)" : "inherit",
                        }}
                        name={
                            selected
                                ? "radio_button_checked"
                                : "radio_button_unchecked"
                        }
                    />
                </Fragment>
            )}
        </Fragment>
    );

    const properties: any = {
        className: classNames({
            [styles.button]: true,
            [styles.primary]: primary,
            [styles.large]: large,
            [styles.selected]: selected,
            [styles.disabled]: disabled,
            [styles.loading]: loading,
            [styles.colorRed]: color === "red",
            [className]: true,
        }),
    };

    if (to) {
        return (
            <Link to={to} {...properties} {...others}>
                {content}
            </Link>
        );
    }

    return (
        <button
            {...properties}
            type="button"
            onClick={disabled || loading ? () => {} : onClick}
            {...others}
        >
            {content}
        </button>
    );
}
