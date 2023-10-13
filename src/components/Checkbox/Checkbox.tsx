import styles from "./Checkbox.module.sass";
import Icon from "../Icon/Icon";
import classNames from "classnames";

type Props = {
    checked?: boolean;
};

export default function Checkbox(props: Readonly<Props>) {
    const { checked } = props;

    return (
        <div
            className={classNames({
                [styles.checkbox]: true,
                [styles.checkboxChecked]: checked,
            })}
        >
            <input type="checkbox" />
            <Icon className={styles.icon} name="check" />
        </div>
    );
}
