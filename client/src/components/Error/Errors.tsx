import { HTMLProps } from "react";
import classNames from "classnames";
import styles from "./Error.module.sass";

export function Errors(props: HTMLProps<HTMLDivElement>) {
    const { children, className, ...others } = props;

    if (!children) return null;

    return (
        <div
            className={classNames(styles.errors, className)}
            {...others}
            children={children}
        />
    );
}
