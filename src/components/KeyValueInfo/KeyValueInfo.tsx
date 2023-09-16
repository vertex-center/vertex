import { PropsWithChildren } from "react";

import styles from "./KeyValueInfo.module.sass";
import Spacer from "../Spacer/Spacer";

export function KeyValueGroup(props: PropsWithChildren) {
    const { children } = props;
    return <div className={styles.group}>{children}</div>;
}

type Type = "code";

type Props = PropsWithChildren<{
    name: string;
    type?: Type;
}>;

export function KeyValueInfo(props: Props) {
    const { name, type, children } = props;

    let content = children;
    if (type === "code") {
        content = <code className={styles.code}>{content}</code>;
    }

    return (
        <div className={styles.info}>
            <div className={styles.name}>{name}</div>
            <Spacer />
            {content}
        </div>
    );
}
