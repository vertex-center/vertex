import React, { PropsWithChildren } from "react";

import styles from "./KeyValueInfo.module.sass";
import Spacer from "../Spacer/Spacer";
import LoadingValue from "../LoadingValue/LoadingValue";

export function KeyValueGroup(props: Readonly<PropsWithChildren>) {
    const { children } = props;
    return <div className={styles.group}>{children}</div>;
}

type Type = "code";

type Props = PropsWithChildren<{
    name: string;
    type?: Type;
    icon?: React.JSX.Element;
    loading?: boolean;
}>;

export function KeyValueInfo(props: Readonly<Props>) {
    const { name, type, icon, loading, children } = props;

    let content = children;
    if (type === "code") {
        content = <code className={styles.code}>{content}</code>;
    }

    return (
        <div className={styles.info}>
            <div className={styles.key}>
                {icon}
                <div className={styles.name}>{name}</div>
            </div>
            <Spacer />
            {loading ? <LoadingValue /> : content}
        </div>
    );
}
