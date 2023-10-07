import classNames from "classnames";

import styles from "./Icon.module.sass";
import { HTMLProps } from "react";

export type IconProps = Omit<HTMLProps<HTMLSpanElement>, "name"> & {
    name: string | JSX.Element;
    rotating?: boolean;
};

export default function Icon(props: Readonly<IconProps>) {
    const { name, rotating, className, ...others } = props;
    return (
        <span
            className={classNames({
                "material-symbols-rounded": true,
                [styles.rotating]: rotating,
                [className]: true,
            })}
            {...others}
        >
            {name}
        </span>
    );
}
