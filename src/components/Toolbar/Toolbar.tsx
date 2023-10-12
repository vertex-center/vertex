import { HTMLProps } from "react";
import styles from "./Toolbar.module.sass";
import Button, { ButtonProps } from "../Button/Button";

function ToolbarButton(props: Readonly<ButtonProps>) {
    return <Button className={styles.button} {...props} />;
}

type Props = HTMLProps<HTMLDivElement>;

function Toolbar(props: Readonly<Props>) {
    return <div className={styles.toolbar} {...props} />;
}

Toolbar.Button = ToolbarButton;

export default Toolbar;
