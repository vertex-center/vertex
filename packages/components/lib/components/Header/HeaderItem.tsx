import cx from "classnames";
import { HTMLProps, ReactNode, useState } from "react";
import { Dropdown } from "../Dropdown/Dropdown.tsx";

export type HeaderItemProps = HTMLProps<HTMLDivElement> & {
    items?: ReactNode;
};

export function HeaderItem(props: Readonly<HeaderItemProps>) {
    const { className, children, items, ...others } = props;

    const [opened, setOpened] = useState(false);

    const onClick = () => setOpened((o: boolean) => !o);
    const close = () => setOpened(false);

    return (
        <div
            className={cx(
                "header-item",
                {
                    "header-item-opened": opened,
                },
                className,
            )}
            onClick={onClick}
            {...others}
        >
            {children}
            {items !== undefined && (
                <Dropdown opened={opened} onClose={close}>
                    {items}
                </Dropdown>
            )}
        </div>
    );
}
