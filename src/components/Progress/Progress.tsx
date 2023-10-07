import styles from "./Progress.module.sass";
import classNames from "classnames";

export function ProgressOverlay({ show }: { show?: boolean }) {
    if (!show) return null;
    return (
        <div className={styles.top}>
            <Progress infinite small />
        </div>
    );
}

type Props = {
    infinite?: boolean;
    small?: boolean;
    value?: number;
};

export default function Progress(props: Readonly<Props>) {
    const { infinite, small, value } = props;

    return (
        <div
            className={classNames({
                [styles.progress]: true,
                [styles.progressSmall]: small,
            })}
        >
            <div
                className={classNames({
                    [styles.bar]: true,
                    [styles.barInfinite]: infinite,
                    [styles.barSmall]: small,
                })}
                style={{
                    width: `${value}%`,
                }}
            />
        </div>
    );
}
