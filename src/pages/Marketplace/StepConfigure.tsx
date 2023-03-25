import { Env, Instance, saveInstanceEnv } from "../../backend/backend";
import { Vertical } from "../../components/Layouts/Layouts";
import { Title } from "../../components/Text/Text";
import styles from "./Marketplace.module.sass";
import { useEffect, useState } from "react";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";
import Button from "../../components/Button/Button";

type StepConfigureProps = {
    onNextStep: () => void;
    instance: Instance;
};

export default function StepConfigure(props: StepConfigureProps) {
    const { onNextStep, instance } = props;

    const [env, setEnv] = useState<any[]>();

    const [uploading, setUploading] = useState(false);

    useEffect(() => {
        setEnv(
            instance.environment.map((e) => ({
                env: e,
                value: e.default ?? "",
            }))
        );
    }, [instance.environment]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
    };

    const save = () => {
        const _env: Env = {};
        env.forEach((e) => {
            _env[e.env.name] = e.value;
        });
        setUploading(true);
        saveInstanceEnv(instance.uuid, _env)
            .then(() => {
                onNextStep();
            })
            .catch(console.error)
            .finally(() => {
                setUploading(false);
            });
    };

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Configure</Title>
            </div>
            <Vertical gap={30}>
                {env?.map((e, i) => (
                    <EnvVariableInput
                        key={i}
                        env={e.env}
                        value={e.value}
                        onChange={(v: any) => onChange(i, v)}
                        disabled={uploading}
                    />
                ))}
            </Vertical>
            <Button
                primary
                large
                rightSymbol="save"
                onClick={save}
                loading={uploading}
            >
                Save
            </Button>
        </div>
    );
}
