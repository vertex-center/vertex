import styles from "./Logo.module.sass";
import classNames from "classnames";

type Props = {
    name?: string;
    iconOnly?: boolean;
    devtools?: boolean;
};

export default function Logo({ name, iconOnly, devtools }: Readonly<Props>) {
    return (
        <div
            className={classNames({
                [styles.logo]: true,
                [styles.logoDev]: devtools,
            })}
        >
            <img width={30} src="/images/logo.png" alt="Logo" />
            {!iconOnly && <h1>{name ?? "Vertex"}</h1>}
            {devtools && <code className={styles.tag}>DEVTOOLS</code>}
        </div>
    );
}
