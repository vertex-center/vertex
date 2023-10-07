import classNames from "classnames";

import styles from "./Symbol.module.sass";
import { HTMLProps } from "react";

export type SymbolProps = Omit<HTMLProps<HTMLSpanElement>, "name"> & {
    name: string | JSX.Element;
    rotating?: boolean;
};

export default function Symbol(props: Readonly<SymbolProps>) {
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
