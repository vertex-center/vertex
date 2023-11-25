import cx from "classnames";
import "./Dropdown.sass";
import { Fragment, HTMLProps } from "react";
import { Overlay } from "../Overlay/Overlay.tsx";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon.tsx";

export type DropdownProps = HTMLProps<HTMLDivElement> & {
    opened?: boolean;
    onClose?: () => void;
};

export function Dropdown(props: Readonly<DropdownProps>) {
    const { className, opened, onClose, ...others } = props;

    return (
        <Fragment>
            <div
                className={cx(
                    "dropdown",
                    {
                        "dropdown-opened": opened,
                    },
                    className,
                )}
                {...others}
            />
            <Overlay show={opened} onClick={onClose} />
        </Fragment>
    );
}

export type DropdownItemProps = HTMLProps<HTMLDivElement> & {
    icon?: string;
    red?: boolean;
};

export function DropdownItem(props: Readonly<DropdownItemProps>) {
    const { className, children, icon, red, ...others } = props;

    return (
        <div
            className={cx(
                "dropdown-item",
                {
                    "dropdown-item-red": red,
                },
                className,
            )}
            {...others}
        >
            {icon && <MaterialIcon icon={icon} />}
            {children}
        </div>
    );
}
