import { Horizontal } from "../Layouts/Layouts";
import Symbol from "../Symbol/Symbol";

import styles from "./Error.module.sass";

type Props = {
    error?: string;
};

export function Error(props: Props) {
    const { error } = props;

    if (!error) return null;

    return (
        <Horizontal gap={8} className={styles.error}>
            <Symbol name="error" />
            <div className={styles.content}>{error}</div>
        </Horizontal>
    );
}
