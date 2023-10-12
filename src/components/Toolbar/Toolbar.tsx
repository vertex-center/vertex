import { HTMLProps } from "react";
import styles from "./Toolbar.module.sass";

type Props = HTMLProps<HTMLDivElement>;

export default function Toolbar(props: Readonly<Props>) {
    return <div className={styles.toolbar} {...props} />;
}
