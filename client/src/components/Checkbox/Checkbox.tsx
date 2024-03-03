import styles from "./Checkbox.module.sass";
import classNames from "classnames";
import { MaterialIcon } from "@vertex-center/components";

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
            <MaterialIcon icon="check" className={styles.icon} />
        </div>
    );
}
