import { Horizontal } from "../Layouts/Layouts";

import styles from "./SegmentedButton.module.sass";
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

type Props = {
    disabled?: boolean;
    value: any;
    onChange: (value: any) => void;
    items: {
        label?: string;
        value: any;
        rightSymbol?: string;
    }[];
};

export function SegmentedButtons(props: Props) {
    const { value, onChange, items, disabled } = props;

    return (
        <Horizontal className={styles.wrapper}>
            {items.map((item) => {
                return (
                    <SegmentedButton
                        key={item.value}
                        value={item.value}
                        onClick={() => onChange(item.value)}
                        selected={value === item.value}
                        disabled={disabled}
                        rightSymbol={item.rightSymbol}
                    >
                        {item.label ?? item.value}
                    </SegmentedButton>
                );
            })}
        </Horizontal>
    );
}
