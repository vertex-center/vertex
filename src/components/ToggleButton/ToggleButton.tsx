import styles from "./ToggleButton.module.sass";
import classNames from "classnames";
import Symbol from "../Symbol/Symbol";

type Props = {
    value: boolean;
    onChange: (value: boolean) => void;
};

export default function ToggleButton(props: Props) {
    const { value } = props;

    const onChange = () => props.onChange(!value);

    return (
        <div
            className={classNames({
                [styles.toggle]: true,
                [styles.on]: value,
            })}
            onClick={onChange}
        >
            <div className={styles.dot}>
                <Symbol
                    className={styles.symbol}
                    name={value ? "check" : "close"}
                />
            </div>
        </div>
    );
}
