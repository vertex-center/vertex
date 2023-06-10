import styles from "./Progress.module.sass";
import classNames from "classnames";

type Props = {
    infinite?: boolean;
    small?: boolean;
};

export default function Progress(props: Props) {
    const { infinite, small } = props;

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
            />
        </div>
    );
}
