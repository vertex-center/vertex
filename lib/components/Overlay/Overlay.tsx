import cx from "classnames";
import React, { HTMLProps } from "react";
import "./Overlay.sass";
import { createPortal } from "react-dom";

export type OverlayProps = HTMLProps<HTMLDivElement> & {
    show?: boolean;
};

export function Overlay(props: Readonly<OverlayProps>) {
    const { className, show, onClick, ...others } = props;
    if (!show) return null;

    const app = document.getElementById("app");
    if (app === null) return null;

    const close = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
        onClick?.(e);
        e.stopPropagation();
    };

    return createPortal(
        <div
            className={cx("overlay", className)}
            onClick={close}
            {...others}
        />,
        app,
    );
}
