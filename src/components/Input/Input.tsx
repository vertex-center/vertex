import { forwardRef, HTMLProps, Ref } from "react";

import styles from "./Input.module.sass";
import classNames from "classnames";
import { Vertical } from "../Layouts/Layouts";
import { Caption } from "../Text/Text";

export type InputProps = HTMLProps<HTMLInputElement> & {
    description?: string;
    error?: string;
};

export type SelectProps = HTMLProps<HTMLSelectElement> & {
    description?: string;
    error?: string;
    required?: boolean;
};

export default forwardRef(function Input(
    props: Readonly<InputProps>,
    ref: Ref<HTMLInputElement>
) {
    const { className, label, description, error, required, ...others } = props;

    return (
        <Vertical gap={10}>
            {label && (
                <label className={styles.label}>
                    {label}
                    {required && <span className={styles.required}>*</span>}
                    {!required && (
                        <span className={styles.optional}>(optional)</span>
                    )}
                </label>
            )}
            <Vertical gap={5}>
                <input
                    ref={ref}
                    className={classNames(styles.input, className)}
                    {...others}
                />
                {description && !error && (
                    <Caption className={styles.description}>
                        {description}
                    </Caption>
                )}
                {error && <Caption className={styles.error}>{error}</Caption>}
            </Vertical>
        </Vertical>
    );
});
