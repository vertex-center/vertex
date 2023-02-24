import classNames from "classnames";

import styles from "./Symbol.module.sass";

type Props = {
    name: string;
    rotating?: boolean;
};

export default function Symbol({ name, rotating }: Props) {
    return (
        <span
            className={classNames({
                "material-symbols-rounded": true,
                [styles.rotating]: rotating,
            })}
        >
            {name}
        </span>
    );
}
