import { PropsWithChildren } from "react";
import Symbol from "../Symbol/Symbol";

import styles from "./URL.module.sass";

type Props = PropsWithChildren<{
    href: string;
}>;

export default function URL({ children, href }: Readonly<Props>) {
    return (
        <a href={href} className={styles.url}>
            <Symbol name="link" />
            {children}
        </a>
    );
}
