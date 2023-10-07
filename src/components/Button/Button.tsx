import { Fragment, HTMLProps } from "react";

import styles from "./Button.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import Spacer from "../Spacer/Spacer";
import { Link } from "react-router-dom";

export type ButtonProps = HTMLProps<HTMLButtonElement> & {
    leftSymbol?: string | any;
    rightSymbol?: string | any;

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
        leftSymbol,
        rightSymbol,
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
            {leftSymbol && (
                <Symbol className={styles.symbol} name={leftSymbol} />
            )}
            <div>{children}</div>
            {rightSymbol && (
                <Symbol className={styles.symbol} name={rightSymbol} />
            )}
            {selectable && (
                <Fragment>
                    <Spacer />
                    <Symbol
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
