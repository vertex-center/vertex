import { HTMLProps } from "react";
import styles from "./Toolbar.module.sass";
import classNames from "classnames";

type Props = HTMLProps<HTMLDivElement>;

export default function Toolbar(props: Readonly<Props>) {
    const { className, ...others } = props;
    return (
        <div className={classNames(styles.toolbar, className)} {...others} />
    );
}
