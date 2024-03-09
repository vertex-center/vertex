import { CircleNotch } from "@phosphor-icons/react";
import "./Spinner.sass";
import { ReactNode } from "react";
import { Horizontal } from "@vertex-center/components/lib/components/Layout/Layout.tsx";

export type SpinnerProps = {
    label?: ReactNode;
};

export function Spinner(props: SpinnerProps) {
    const { label } = props;

    const spinner = <CircleNotch className="spinner" />;

    if (label) {
        return (
            <Horizontal gap={8} alignItems="center">
                {spinner}
                <div className="spinner-label">{label}</div>
            </Horizontal>
        );
    }

    return spinner;
}
