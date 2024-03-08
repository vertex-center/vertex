import styles from "./Popup.module.sass";
import { HTMLProps, PropsWithChildren, useEffect } from "react";
import classNames from "classnames";
import { Title, Vertical } from "@vertex-center/components";

export function PopupActions(props: HTMLProps<HTMLDivElement>) {
    const { className, ...others } = props;
    return (
        <div className={classNames(styles.actions, className)} {...others} />
    );
}

type Props = PropsWithChildren<{
    show: boolean;
    onDismiss: () => void;
    title: string;
}>;

export default function Popup(props: Readonly<Props>) {
    const { show, onDismiss, title, children } = props;

    if (!show) return null;

    useEffect(() => {
        const handleKeyPress = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                onDismiss();
            }
        };
        window.addEventListener("keydown", handleKeyPress);
        return () => window.removeEventListener("keydown", handleKeyPress);
    }, [show, onDismiss]);

    return (
        <div>
            <div className={styles.overlay} onClick={onDismiss} />
            <div className={styles.popup}>
                <div className={styles.header}>
                    <Title variant="h3">{title}</Title>
                </div>
                <Vertical gap={20}>{children}</Vertical>
            </div>
        </div>
    );
}
