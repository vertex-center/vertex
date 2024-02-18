import { MaterialIcon } from "../MaterialIcon/MaterialIcon";
import cx from "classnames";
import "./Checkbox.sass";
import { HTMLAttributes } from "react";

export type CheckboxProps = HTMLAttributes<HTMLDivElement> & {
    checked?: boolean;
};

export function Checkbox(props: Readonly<CheckboxProps>) {
    const { checked, className, ...others } = props;

    return (
        <div
            className={cx(
                "checkbox",
                {
                    "checkbox-checked": checked,
                },
                className,
            )}
            {...others}
        >
            <input type="checkbox" />
            <MaterialIcon icon="check" className="checkbox-icon" />
        </div>
    );
}
