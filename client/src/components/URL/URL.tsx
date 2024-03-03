import { HTMLProps } from "react";

import styles from "./URL.module.sass";
import classNames from "classnames";
import { LinkSimple } from "@phosphor-icons/react";

type Props = HTMLProps<HTMLAnchorElement>;

export default function URL(props: Readonly<Props>) {
    const { className, children, ...others } = props;
    return (
        <a className={classNames(styles.url, className)} {...others}>
            <LinkSimple />
            {children}
        </a>
    );
}
