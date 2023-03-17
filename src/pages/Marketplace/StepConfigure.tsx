import { EnvVariable, Service } from "../../backend/backend";
import PortInput from "../../components/Input/PortInput";
import Input from "../../components/Input/Input";
import { Vertical } from "../../components/Layouts/Layouts";
import { Caption, Title } from "../../components/Text/Text";
import styles from "./Marketplace.module.sass";
import { useEffect, useState } from "react";
import Symbol from "../../components/Symbol/Symbol";

type VariableInputProps = {
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
};

function VariableInput(props: VariableInputProps) {
    const { env, value, onChange } = props;

    const inputProps = {
        value,
        label: env.display_name,
        name: env.name,
        onChange: (e) => onChange(e.target.value),
    };

    let input;
    if (env.type === "port") {
        input = <PortInput {...inputProps} />;
    } else {
        input = <Input {...inputProps} />;
    }

    return (
        <Vertical gap={6}>
            {input}
            <Caption className={styles.inputDescription}>
                {env.description}
            </Caption>
        </Vertical>
    );
}

type StepConfigureProps = {
    service: Service;
};

export default function StepConfigure(props: StepConfigureProps) {
    const { service } = props;

    const [env, setEnv] = useState<any[]>();

    useEffect(() => {
        setEnv(
            service.environment.map((e) => ({
                env: e,
                value: e.default ?? "",
            }))
        );
    }, [service.environment]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
    };

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Symbol name="counter_2" />
                <Title>Configure</Title>
            </div>
            <Vertical gap={30}>
                {env?.map((e, i) => (
                    <VariableInput
                        key={i}
                        env={e.env}
                        value={e.value}
                        onChange={(v: any) => onChange(i, v)}
                    />
                ))}
            </Vertical>
        </div>
    );
}
