import styles from "./ProxyRedirect.module.sass";
import classNames from "classnames";
import { MaterialIcon } from "@vertex-center/components";

type Props = {
    source: string;
    target: string;
    enabled?: boolean;
    onDelete?: () => void;
};

export default function ProxyRedirect(props: Readonly<Props>) {
    const { source, target, enabled, onDelete } = props;

    return (
        <div
            className={classNames({
                [styles.redirect]: true,
                [styles.enabled]: enabled,
            })}
        >
            <div className={styles.wrapper}>
                <div className={styles.url}>
                    <MaterialIcon icon="link" />
                    {source}
                </div>
                <div className={styles.line} />
                <MaterialIcon icon="double_arrow" className={styles.arrow} />
                <div className={styles.line} />
                <div className={styles.url}>
                    <MaterialIcon icon="subdirectory_arrow_right" />
                    {target}
                </div>
            </div>
            <div className={styles.delete} onClick={onDelete}>
                <MaterialIcon icon="close" />
            </div>
        </div>
    );
}
