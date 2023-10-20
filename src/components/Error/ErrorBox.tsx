import styles from "./Error.module.sass";
import { HTMLProps } from "react";
import classNames from "classnames";
import { MaterialIcon } from "@vertex-center/components";

type Props = HTMLProps<HTMLDivElement> & {
    error?: any;
};

export default function ErrorBox(props: Readonly<Props>) {
    const { error, className, ...others } = props;

    let err = error?.message ?? "An unknown error has occurred.";

    return (
        <div className={classNames(styles.box, className)} {...others}>
            <div className={styles.error}>
                <MaterialIcon icon="error" className={styles.icon} />
                <h1>Error</h1>
            </div>
            <div className={styles.content}>{err}</div>
        </div>
    );
}
