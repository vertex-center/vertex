import styles from "./ToggleButton.module.sass";
import classNames from "classnames";
import { MaterialIcon } from "@vertex-center/components";

type Props = {
    value: boolean;
    onChange: (value: boolean) => void;
    disabled?: boolean;
};

export default function ToggleButton(props: Readonly<Props>) {
    const { value, disabled } = props;

    const onChange = () => props.onChange(!value);

    return (
        <div
            className={classNames({
                [styles.toggle]: true,
                [styles.disabled]: disabled,
                [styles.on]: value,
            })}
            onClick={onChange}
        >
            <div className={styles.dot}>
                <MaterialIcon
                    icon={value ? "check" : "close"}
                    className={styles.icon}
                />
            </div>
        </div>
    );
}
