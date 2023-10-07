import { HTMLProps } from "react";
import styles from "./VersionTag.module.sass";
import classNames from "classnames";

type Props = HTMLProps<HTMLSpanElement>;

export default function VersionTag(props: Readonly<Props>) {
    const { children, className, ...others } = props;

    return (
        <span className={classNames(styles.tag, className)} {...others}>
            {children}
        </span>
    );
}
