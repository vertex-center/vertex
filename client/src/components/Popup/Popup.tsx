import styles from "./Popup.module.sass";
import { HTMLProps, ReactNode, useEffect } from "react";
import classNames from "classnames";
import { Title, Vertical } from "@vertex-center/components";

export function PopupActions(props: HTMLProps<HTMLDivElement>) {
    const { className, ...others } = props;
    return (
        <div className={classNames(styles.actions, className)} {...others} />
    );
}

type Props = HTMLProps<HTMLDivElement> & {
    onDismiss: () => void;
    title: string;
    image?: ReactNode;
};

export default function Popup(props: Readonly<Props>) {
    const { onDismiss, title, image, children, ...others } = props;

    useEffect(() => {
        const handleKeyPress = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                onDismiss();
            }
        };
        window.addEventListener("keydown", handleKeyPress);
        return () => window.removeEventListener("keydown", handleKeyPress);
    }, [onDismiss]);

    return (
        <div>
            <div className={styles.overlay} onClick={onDismiss} />
            <div className={styles.popup} {...others}>
                <div className={styles.image}>{image}</div>
                <div className={styles.content}>
                    <div className={styles.header}>
                        <Title variant="h3">{title}</Title>
                    </div>
                    <Vertical gap={20}>{children}</Vertical>
                </div>
            </div>
        </div>
    );
}
