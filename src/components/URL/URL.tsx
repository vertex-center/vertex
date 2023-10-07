import { PropsWithChildren } from "react";
import Icon from "../Icon/Icon";

import styles from "./URL.module.sass";

type Props = PropsWithChildren<{
    href: string;
}>;

export default function URL({ children, href }: Readonly<Props>) {
    return (
        <a href={href} className={styles.url}>
            <Icon name="link" />
            {children}
        </a>
    );
}
