import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { saveInstanceEnv } from "../../backend/backend";
import { useParams } from "react-router-dom";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";
import Button from "../../components/Button/Button";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";
import useInstance from "../../hooks/useInstance";
import { Env, EnvVariable } from "../../models/service";
import styles from "./BayDetailsEnv.module.sass";
import Loading from "../../components/Loading/Loading";

type Props = {};

export default function BayDetailsEnv(props: Props) {
    const { uuid } = useParams();

    const [env, setEnv] = useState<{ env: EnvVariable; value: any }[]>();

    const { instance } = useInstance(uuid);

    const [uploading, setUploading] = useState(false);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);

    useEffect(() => {
        setEnv(
            instance?.environment?.map((e) => ({
                env: e,
                value: instance?.env[e.name] ?? e.default ?? "",
            }))
        );
    }, [instance]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
        setSaved(false);
    };

    const save = () => {
        const _env: Env = {};
        env.forEach((e) => {
            _env[e.env.name] = e.value;
        });
        setUploading(true);
        saveInstanceEnv(uuid, _env)
            .then(console.log)
            .catch(console.error)
            .finally(() => {
                setUploading(false);
                setSaved(true);
            });
    };

    return (
        <Fragment>
            <Title className={styles.title}>Environment</Title>
            {env?.map((env, i) => (
                <EnvVariableInput
                    env={env.env}
                    value={env.value}
                    onChange={(v) => onChange(i, v)}
                    disabled={uploading}
                />
            ))}
            <Button
                primary
                large
                onClick={save}
                rightSymbol="save"
                loading={uploading}
                disabled={saved || saved === undefined}
            >
                Save{" "}
                {instance?.install_method === "docker" &&
                    "+ Recreate container"}
            </Button>
            {uploading && <Loading />}
            {saved && (
                <Horizontal
                    className={styles.saved}
                    alignItems="center"
                    gap={4}
                >
                    <Symbol name="check" />
                    Saved!
                </Horizontal>
            )}
        </Fragment>
    );
}
