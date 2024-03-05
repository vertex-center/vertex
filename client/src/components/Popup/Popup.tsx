import styles from "./Popup.module.sass";
import { HTMLProps, PropsWithChildren, ReactNode } from "react";
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

    return (
        <div>
            <div
                className={classNames({
                    [styles.overlay]: true,
                    [styles.overlayActive]: show,
                })}
                onClick={onDismiss}
            />
            <div
                className={classNames({
                    [styles.popup]: true,
                    [styles.popupActive]: show,
                })}
            >
                <div className={styles.header}>
                    <Title variant="h3">{title}</Title>
                </div>
                <Vertical gap={20}>{children}</Vertical>
            </div>
        </div>
    );
}
