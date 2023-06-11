import { HTMLProps } from "react";

import styles from "./Text.module.sass";
import classNames from "classnames";

export function Title(props: HTMLProps<HTMLHeadingElement>) {
    const { children, className, ...others } = props;

    return (
        <h2 className={classNames(styles.title, className)} {...others}>
            {children}
        </h2>
    );
}

export function Caption(props: HTMLProps<HTMLParagraphElement>) {
    const { children, className, ...others } = props;

    return (
        <p className={classNames(styles.caption, className)} {...others}>
            {children}
        </p>
    );
}

export function Text(props: HTMLProps<HTMLParagraphElement>) {
    const { children, className, ...others } = props;

    return (
        <p className={classNames(styles.text, className)} {...others}>
            {children}
        </p>
    );
}
