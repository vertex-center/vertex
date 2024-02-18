import React, { Children, cloneElement, HTMLProps, useId } from "react";
import "./FormItem.sass";
import cx from "classnames";

export type FormItemProps = HTMLProps<HTMLDivElement> & {
    label?: string;
    required?: boolean;
    optional?: boolean;
    error?: string;
    description?: string;
};

export function FormItem(props: Readonly<FormItemProps>) {
    const {
        label,
        required,
        optional,
        error,
        description,
        children,
        className,
        ...others
    } = props;

    const id = useId();

    let indicator;
    if (required) {
        indicator = <span className="form-item-required">*</span>;
    } else if (optional) {
        indicator = <span className="form-item-optional">(optional)</span>;
    }

    return (
        <div className={cx("form-item", className)} {...others}>
            {label && (
                <label htmlFor={id} className="form-item-label">
                    {label} {indicator}
                </label>
            )}
            {Children.map(children, (child: React.ReactNode) => {
                if (!React.isValidElement(child)) return null;
                // @ts-expect-error cloneElement is too hard to type
                return cloneElement<never>(child, {
                    id,
                    required,
                });
            })}
            {description && !error && (
                <div className="form-item-description">{description}</div>
            )}
            {error && <div className="form-item-error">{error}</div>}
        </div>
    );
}
