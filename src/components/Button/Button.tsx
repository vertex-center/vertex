import { Fragment, HTMLProps } from "react";

import styles from "./Button.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import Spacer from "../Spacer/Spacer";

export type ButtonProps = HTMLProps<HTMLButtonElement> & {
    leftSymbol?: string | any;
    rightSymbol?: string | any;

    selected?: boolean;
    selectable?: boolean;

    loading?: boolean;
    disabled?: boolean;

    onClick: () => void;

    // types
    primary?: boolean;
    large?: boolean;
};

export default function Button(props: ButtonProps) {
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
        className,
        type,
        ...others
    } = props;

    return (
        <button
            className={classNames({
                [styles.button]: true,
                [styles.primary]: primary,
                [styles.large]: large,
                [styles.selected]: selected,
                [styles.disabled]: disabled,
                [styles.loading]: loading,
                [className]: true,
            })}
            type="button"
            onClick={disabled || loading ? () => {} : onClick}
            {...others}
        >
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
        </button>
    );
}
