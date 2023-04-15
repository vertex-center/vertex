import styles from "./Popup.module.sass";
import { Fragment, PropsWithChildren } from "react";
import classNames from "classnames";

type Props = PropsWithChildren<{
    show: boolean;
    onDismiss: () => void;
}>;

export default function Popup(props: Props) {
    const { show, onDismiss, children } = props;

    return (
        <Fragment>
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
                {children}
            </div>
        </Fragment>
    );
}
