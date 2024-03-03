import styles from "./Popup.module.sass";
import { PropsWithChildren, ReactNode } from "react";
import classNames from "classnames";
import { Title, Vertical } from "@vertex-center/components";

type Props = PropsWithChildren<{
    show: boolean;
    onDismiss: () => void;
    title: string;
    actions: ReactNode;
}>;

export default function Popup(props: Readonly<Props>) {
    const { show, onDismiss, title, actions, children } = props;

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
                <Vertical gap={24}>{children}</Vertical>
                <div className={styles.actions}>{actions}</div>
            </div>
        </div>
    );
}
