import { HTMLProps } from "react";
import styles from "./Text.module.sass";
import cx from "classnames";

export function Caption(props: HTMLProps<HTMLParagraphElement>) {
    const { children, className, ...others } = props;

    return (
        <p className={cx(styles.caption, className)} {...others}>
            {children}
        </p>
    );
}
