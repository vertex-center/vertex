import { HTMLProps } from "react";

import styles from "./Text.module.sass";
import classNames from "classnames";

export function BigTitle(props: HTMLProps<HTMLHeadingElement>) {
    const { children, className, ...others } = props;

    return (
        <h1 className={classNames(styles.bigtitle, className)} {...others}>
            {children}
        </h1>
    );
}

export function Title(props: HTMLProps<HTMLHeadingElement>) {
    const { children, className, ...others } = props;

    return (
        <h2 className={classNames(styles.title, className)} {...others}>
            {children}
        </h2>
    );
}

export function SubTitle(props: HTMLProps<HTMLHeadingElement>) {
    const { children, className, ...others } = props;

    return (
        <h3 className={classNames(styles.subtitle, className)} {...others}>
            {children}
        </h3>
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
