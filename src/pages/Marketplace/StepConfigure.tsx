import { Service } from "../../backend/backend";
import { Vertical } from "../../components/Layouts/Layouts";
import { Title } from "../../components/Text/Text";
import styles from "./Marketplace.module.sass";
import { useEffect, useState } from "react";
import Symbol from "../../components/Symbol/Symbol";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";

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
                    <EnvVariableInput
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
