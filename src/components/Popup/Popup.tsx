import styles from "./Popup.module.sass";
import { PropsWithChildren } from "react";
import classNames from "classnames";

type Props = PropsWithChildren<{
    show: boolean;
    onDismiss: () => void;
}>;

export default function Popup(props: Readonly<Props>) {
    const { show, onDismiss, children } = props;

    return (
        <div
            className={classNames({
                [styles.overlay]: true,
                [styles.overlayActive]: show,
            })}
            onClick={onDismiss}
        >
            <div
                className={classNames({
                    [styles.popup]: true,
                    [styles.popupActive]: show,
                })}
            >
                {children}
            </div>
        </div>
    );
}
