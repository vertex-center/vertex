import { Horizontal } from "../Layouts/Layouts";

import styles from "./SegmentedButton.module.sass";
import { HTMLProps } from "react";
import Button, { ButtonProps } from "../Button/Button";
import classNames from "classnames";

type SegmentedButtonProps = ButtonProps & {
    value: any;
};

export function SegmentedButton(props: SegmentedButtonProps) {
    const { children, className, ...others } = props;

    return (
        <Button {...others} className={classNames(styles.button, className)}>
            {children}
        </Button>
    );
}

type Props = HTMLProps<HTMLDivElement> & {
    value: any;
    onChange: (value: any) => void;
    items: {
        value: string;
    }[];
};

export function SegmentedButtons(props: Props) {
    const { value, onChange, items, disabled } = props;

    return (
        <Horizontal className={styles.wrapper}>
            {items.map((item) => {
                return (
                    <SegmentedButton
                        value={item.value}
                        onClick={() => onChange(item.value)}
                        selected={value === item.value}
                        disabled={disabled}
                    >
                        {item.value}
                    </SegmentedButton>
                );
            })}
        </Horizontal>
    );
}
