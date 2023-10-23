import { Input, InputProps } from "../Input/Input.tsx";
import { Children, cloneElement, HTMLProps, useState } from "react";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon.tsx";
import "./SelectField.sass";
import cx from "classnames";

type SelectValueProps = HTMLProps<HTMLDivElement>;

export function SelectValue(props: Readonly<SelectValueProps>) {
    return <div className="select-field-value" {...props} />;
}

type SelectOptionProps = HTMLProps<HTMLDivElement>;

export function SelectOption(props: Readonly<SelectOptionProps>) {
    return <div className="select-field-option" {...props} />;
}

type SelectFieldProps = InputProps & {
    onChange?: (value: any) => void;
};

export function SelectField(props: Readonly<SelectFieldProps>) {
    const { children, ...others } = props;

    const [show, setShow] = useState<boolean>(false);

    const onChange = (value: any) => {
        setShow(false);
        props?.onChange?.(value);
    };

    const toggle = () => setShow(!show);

    return (
        <div
            className={cx("select-field", {
                "select-field-opened": show,
            })}
        >
            <Input
                {...others}
                as={"div"}
                inputProps={{
                    className: "select-field-input",
                    onClick: toggle,
                }}
            >
                <MaterialIcon
                    className="select-field-icon"
                    icon="expand_more"
                />
            </Input>
            <div className="select-field-values">
                {Children.map(children, (child) => {
                    if (!child) return;
                    return cloneElement(child, {
                        onClick: onChange,
                    });
                })}
            </div>
            <div className="select-field-overlay" onClick={toggle} />
        </div>
    );
}
