import styles from "./ProxyRedirect.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";

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
                    <Icon name="link" />
                    {source}
                </div>
                <div className={styles.line} />
                <Icon className={styles.arrow} name="double_arrow" />
                <div className={styles.line} />
                <div className={styles.url}>
                    <Icon name="subdirectory_arrow_right" />
                    {target}
                </div>
            </div>
            <div className={styles.delete} onClick={onDelete}>
                <Icon name="close" />
            </div>
        </div>
    );
}
