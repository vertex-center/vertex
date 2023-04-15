import { HTMLProps, PropsWithChildren } from "react";

import styles from "./Text.module.sass";
import classNames from "classnames";

export function Title({ children }: PropsWithChildren) {
    return <h2 className={styles.title}>{children}</h2>;
}

export function Caption(props: HTMLProps<HTMLHeadingElement>) {
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
