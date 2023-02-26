import { Fragment, PropsWithChildren } from "react";

import styles from "./Button.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import Spacer from "../Spacer/Spacer";

type Props = PropsWithChildren<{
    leftSymbol?: string;
    rightSymbol?: string;

    selected?: boolean;
    selectable?: boolean;

    disabled?: boolean;

    onClick: () => void;

    // types
    primary?: boolean;
    large?: boolean;
}>;

export default function Button(props: Props) {
    const {
        children,
        leftSymbol,
        rightSymbol,
        disabled,
        primary,
        large,
        selected,
        selectable,
        onClick,
    } = props;

    return (
        <button
            className={classNames({
                [styles.button]: true,
                [styles.primary]: primary,
                [styles.large]: large,
                [styles.selected]: selected,
                [styles.disabled]: disabled,
            })}
            type="button"
            onClick={disabled ? () => {} : onClick}
        >
            {leftSymbol && <Symbol name={leftSymbol} />}
            <div>{children}</div>
            {rightSymbol && <Symbol name={rightSymbol} />}
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
