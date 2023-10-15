import { HTMLProps } from "react";
import Icon from "../Icon/Icon";

import styles from "./URL.module.sass";
import classNames from "classnames";

type Props = HTMLProps<HTMLAnchorElement>;

export default function URL(props: Readonly<Props>) {
    const { className, children, ...others } = props;
    return (
        <a className={classNames(styles.url, className)} {...props}>
            <Icon name="link" />
            {children}
        </a>
    );
}
