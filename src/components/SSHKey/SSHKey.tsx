import { HTMLProps } from "react";
import styles from "./SSHKey.module.sass";
import Symbol from "../Symbol/Symbol";
import { Vertical } from "../Layouts/Layouts";

export function SSHKeys(props: HTMLProps<HTMLDivElement>) {
    return <div {...props} />;
}

type Props = {};

export default function SSHKey(props: Props) {
    return (
        <div className={styles.key}>
            <Symbol name="key" className={styles.symbol} />
            <Vertical gap={4}>
                <div className={styles.name}>SSH Key</div>
                <div className={styles.text}>Yes</div>
            </Vertical>
        </div>
    );
}
