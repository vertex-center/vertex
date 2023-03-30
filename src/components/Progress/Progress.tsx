import styles from "./Progress.module.sass";
import classNames from "classnames";

type Props = {
    infinite?: boolean;
};

export default function Progress(props: Props) {
    const { infinite } = props;

    return (
        <div className={styles.progress}>
            <div
                className={classNames({
                    [styles.bar]: true,
                    [styles.barInfinite]: infinite,
                })}
            />
        </div>
    );
}
