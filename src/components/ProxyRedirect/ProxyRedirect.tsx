import styles from "./ProxyRedirect.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";

type Props = {
    source: string;
    target: string;
    enabled?: boolean;
    onDelete?: () => void;
};

export default function ProxyRedirect(props: Props) {
    const { source, target, enabled, onDelete } = props;

    return (
        <div
            className={classNames({
                [styles.redirect]: true,
                [styles.enabled]: enabled,
            })}
        >
            <div className={styles.url}>{source}</div>
            <div className={styles.line} />
            <Symbol className={styles.arrow} name="double_arrow" />
            <div className={styles.line} />
            <div className={styles.url}>{target}</div>
            <div className={styles.delete} onClick={onDelete}>
                <Symbol name="close" />
            </div>
        </div>
    );
}
